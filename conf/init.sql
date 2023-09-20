--  This file is part of the eliona project.
--  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
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

CREATE SCHEMA IF NOT EXISTS saml_sp;

GRANT USAGE ON SCHEMA saml_sp TO leicom;
GRANT ALL   ON SCHEMA saml_sp TO leicom;

CREATE TABLE IF NOT EXISTS saml_sp.basic_config (
    enable           BOOLEAN PRIMARY KEY NOT NULL DEFAULT true,
    sp_certificate   TEXT                NOT NULL             ,
    sp_private_key   TEXT                NOT NULL             ,
    idp_metadata_url TEXT                                     ,
    metadata_xml     TEXT                         DEFAULT NULL,
    own_url          TEXT                NOT NULL
) ;

CREATE TABLE IF NOT EXISTS saml_sp.attribute_map (
    enable     BOOLEAN          NOT NULL REFERENCES saml_sp.basic_config(enable),
    email      TEXT PRIMARY KEY NOT NULL DEFAULT 'http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn',
    first_name TEXT                      DEFAULT NULL,
    last_name  TEXT                      DEFAULT NULL,
    phone      TEXT                      DEFAULT NULL
) ;

CREATE TABLE IF NOT EXISTS saml_sp.advanced_config (
    enable                      BOOLEAN PRIMARY KEY NOT NULL REFERENCES saml_sp.basic_config(enable),
    allow_initialization_by_idp BOOLEAN             NOT NULL DEFAULT false,
    signed_request              BOOLEAN             NOT NULL DEFAULT true,
    force_authn                 BOOLEAN             NOT NULL DEFAULT false,
    entity_id                   TEXT                NOT NULL DEFAULT '{ownUrl}/saml/metadata',
    cookie_secure               BOOLEAN             NOT NULL DEFAULT false,
    login_failed_url            TEXT                NOT NULL DEFAULT '{ownUrl}/noLogin'
) ;

CREATE TABLE IF NOT EXISTS saml_sp.permissions (
    enable                     BOOLEAN PRIMARY KEY NOT NULL REFERENCES saml_sp.basic_config(enable),
    default_system_role        TEXT                NOT NULL DEFAULT 'regular', -- reference to is maybe a bad idea (new ACL)
    default_proj_role          TEXT                NOT NULL DEFAULT 'operator',
    system_role_saml_attribute TEXT,
    system_role_map            JSON,
    proj_role_saml_attribute   TEXT,
    proj_role_map              JSON
) ;

-- To INIT
-- INSERT INTO saml_sp.attribute_map (enable, username, username_cut_email_sufix) VALUES (true, 'http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn', false) ON CONFLICT(username) DO NOTHING;
