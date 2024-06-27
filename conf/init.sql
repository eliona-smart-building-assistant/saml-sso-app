--  This file is part of the eliona project.
--  Copyright © 2024 Eliona by IoTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

CREATE SCHEMA IF NOT EXISTS saml_sp ;

-- general settings for the SAML Service Provider (SP)
CREATE TABLE IF NOT EXISTS saml_sp.config (
    id                          INT PRIMARY KEY NOT NULL DEFAULT 1     CHECK (id = 1)                         , -- due to the architecture of eliona only one configuration (1 sso per instance) is possible
    enable                      BOOLEAN         NOT NULL DEFAULT true                                         ,
    sp_certificate              TEXT            NOT NULL                                                      , -- own cert
    sp_private_key              TEXT            NOT NULL                                                      , -- key to own cert
    idp_metadata_url            TEXT                                                                          , -- url where IdP's metadata can fetched
    metadata_xml                TEXT                     DEFAULT NULL                                         , -- if no url is avalable, insert metadata xml here
    own_url                     TEXT            NOT NULL                                                      , -- the own url e.g. https://my.eliona.xy
    user_to_archive             BOOLEAN         NOT NULL DEFAULT false                                        , -- put user to archive @ first login (do not allow login, if not verified by sys admin)
    allow_initialization_by_idp BOOLEAN         NOT NULL DEFAULT false                                        , -- if the IdP can initialize the login (means, no SAML request was issued by our sp)
    signed_request              BOOLEAN         NOT NULL DEFAULT true                                         , -- sign the SAML request
    force_authn                 BOOLEAN         NOT NULL DEFAULT false                                        ,
    entity_id                   TEXT            NOT NULL DEFAULT '{ownUrl}/apps-public/saml-sso/saml/metadata',
    login_failed_url            TEXT            NOT NULL DEFAULT '{ownUrl}/noLogin'                             -- redirect url when a user login fails
) ;

-- general settings for adding a user
CREATE TABLE IF NOT EXISTS saml_sp.attribute_map ( -- SAML session attribute names.
    id              INT PRIMARY KEY NOT NULL DEFAULT 1 REFERENCES saml_sp.config(id) ON UPDATE CASCADE          ,
    email           TEXT            NOT NULL DEFAULT 'http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn', -- SAML attribute email and login
    first_name      TEXT                     DEFAULT NULL                                                       ,
    last_name       TEXT                     DEFAULT NULL                                                       ,
    phone           TEXT                     DEFAULT NULL
) ;

-- settings for define users permissions
CREATE TABLE IF NOT EXISTS saml_sp.permissions (
    id                          INT PRIMARY KEY NOT NULL DEFAULT 1 REFERENCES saml_sp.config(id) ON UPDATE CASCADE,
    default_system_role         TEXT            NOT NULL DEFAULT 'System user'                                    , -- reference to is maybe a bad idea (due to the new ACL)
    default_proj_role           TEXT            NOT NULL DEFAULT 'Project user'                                   , -- can be the role display name or role id
    default_language            TEXT            NOT NULL DEFAULT 'en'                                             , -- see constraint
    system_role_saml_attribute  TEXT                                                                              , -- attribute that contains the system roles which should be mapped
    system_role_map             JSON                                                                              , -- e.g. {"company xy-Admin":"System admin", ...}
    proj_role_saml_attribute    TEXT                                                                              , -- attribute that contains the project roles which should be mapped
    proj_role_map               JSON                                                                              , -- e.g. {"company xy-Employee":"Project user", ...}
    language_saml_attribute     TEXT                                                                              , -- attribute that contains the users language which should be mapped
    language_map                JSON                                                                              , -- e.g. {"Sprache:Deutsch":"de", "Sprache:Englisch":"en"}
    CONSTRAINT chk_language CHECK (default_language IN ('en', 'de', 'it', 'fr'))
) ;

DO $$
BEGIN
    -- Check if the adfs schema exists
    IF EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'adfs') THEN
        -- Migrate data from the existing adfs.config table to the new saml_sp.config table
        INSERT INTO saml_sp.config (id, enable, sp_certificate, sp_private_key, idp_metadata_url, own_url, user_to_archive, allow_initialization_by_idp, signed_request, force_authn, entity_id, login_failed_url)
        SELECT 
            1 AS id, 
            enabled AS enable, 
            cert AS sp_certificate, 
            key AS sp_private_key, 
            metadata_url, 
            own_url, 
            false AS user_to_archive, 
            false AS allow_initialization_by_idp, 
            true AS signed_request, 
            false AS force_authn, 
            own_url || '/apps-public/saml-sso/saml/metadata' AS entity_id, 
            COALESCE(redirect_on_fail_url, own_url || '/noLogin') AS login_failed_url
        FROM adfs.config
        WHERE config_id = 1
        ON CONFLICT (id)
        DO NOTHING;

        -- Migrate data from the existing adfs.attribute_map table to the new saml_sp.attribute_map table
        INSERT INTO saml_sp.attribute_map (id, email, first_name, last_name, phone)
        SELECT 
            1 AS id, 
            email, 
            first_name, 
            last_name, 
            phone
        FROM adfs.attribute_map
        ON CONFLICT (id)
        DO NOTHING;

        -- Migrate additional attribute mappings to the saml_sp.permissions table
        INSERT INTO saml_sp.permissions (id, default_system_role, default_proj_role, default_language, system_role_saml_attribute, system_role_map, proj_role_saml_attribute, proj_role_map, language_saml_attribute, language_map)
        SELECT 
            1 AS id, 
            default_system_role, 
            default_project_role, 
            default_language, 
            system_role_attr, 
            system_role_attr_map, 
            project_role_attr, 
            project_role_attr_map, 
            language_attr, 
            language_attr_map
        FROM adfs.attribute_map
        ON CONFLICT (id)
        DO NOTHING;
    END IF;
END $$;
