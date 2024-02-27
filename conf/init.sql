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

GRANT USAGE ON SCHEMA saml_sp TO leicom ;
GRANT ALL   ON SCHEMA saml_sp TO leicom ;

-- general settings for the SAML Service Provider (SP)
CREATE TABLE IF NOT EXISTS saml_sp.config (
    id                          INT PRIMARY KEY NOT NULL DEFAULT 1     CHECK (id = 1)    , -- due to the architecture of eliona only one configuration (1 sso per instance) is possible
    enable                      BOOLEAN         NOT NULL DEFAULT true                    ,
    sp_certificate              TEXT            NOT NULL                                 , -- own cert
    sp_private_key              TEXT            NOT NULL                                 , -- key to own cert
    idp_metadata_url            TEXT                                                     , -- url where IdP's metadata can fetched
    metadata_xml                TEXT                     DEFAULT NULL                    , -- if no url is avalable, insert metadata xml here
    own_url                     TEXT            NOT NULL                                 , -- the own url e.g. https://my.eliona.xy
    user_to_archive             BOOLEAN         NOT NULL DEFAULT false                   , -- put user to archive @ first login (do not allow login, if not verified by sys admin)
    allow_initialization_by_idp BOOLEAN         NOT NULL DEFAULT false                   , -- if the IdP can initialize the login (means, no SAML request was issued by our sp)
    signed_request              BOOLEAN         NOT NULL DEFAULT true                    , -- sign the SAML request
    force_authn                 BOOLEAN         NOT NULL DEFAULT false                   ,
    entity_id                   TEXT            NOT NULL DEFAULT '{ownUrl}/saml/metadata',
    cookie_secure               BOOLEAN         NOT NULL DEFAULT false                   ,
    login_failed_url            TEXT            NOT NULL DEFAULT '{ownUrl}/noLogin'        -- redirect url when a user login fails
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
    system_role_map             JSON                                                                              , -- e.g. {"firm xy-Admin":"System admin", ...}
    proj_role_saml_attribute    TEXT                                                                              , -- attribute that contains the project roles which should be mapped
    proj_role_map               JSON                                                                              , -- e.g. {"firm xy-Employee":"Project user", ...}
    language_saml_attribute     TEXT                                                                              , -- attribute that contains the users language which should be mapped
    language_map                JSON                                                                              , -- e.g. {"Sprache:Deutsch":"de", "Sprache:Englisch":"en"}
    CONSTRAINT chk_language CHECK (default_language IN ('en', 'de', 'it', 'fr'))
) ;
