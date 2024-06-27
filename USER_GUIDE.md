# SAML SSO

## Eliona app for single-sign-on

> The SAML SSO (Security Assertion Markup Language Single Sign-On) app allows users to log into Eliona using various SSO providers, including Microsoft ADFS. This setup streamlines authentication by using a single set of credentials.

![MS Log-in](user_guide/login.avif)

## Configuration

The SAML 2.0 Service Provider is configured by defining one or more authentication credentials:

| Attribute                   | Description                                                                                             |
|-----------------------------|---------------------------------------------------------------------------------------------------------|
| `id`                        | Configuration Id. Can only be 1                                                                         |
| `enable`                    | If the configuration is enabled or not                                                                  |
| `serviceProviderCertificate`| The Certificate of this SAML Service Provider (SP). Can be a self-signed x509 certificate.              |
| `serviceProviderPrivateKey` | The Private Key matching the Certificate of this SAML Service Provider (SP). DO NOT use RSA key length lower than 2048 |
| `idpMetadataUrl`            | The Metadata URL of the Identity Provider (IdP) if available. Otherwise use the metadataXml to provide Metadata of IdP directly and leave this null |
| `idpMetadataXml`            | Provide the IdP Metadata XML directly, if you do not have the idpMetadataUrl accessible                 |
| `ownUrl`                    | The own URL of this Eliona instance                                                                      |
| `userToArchive`             | If enabled, the newly created user is archived and cannot log in until an admin has activated it         |
| `allowInitializationByIdp`  | If the configuration is enabled or not                                                                  |
| `signedRequest`             | If the SP should make a signed SAML Authn-Request or not                                                |
| `forceAuthn`                | Normally this value is set to false for an SP. If set to true the user has to re-authenticate (login at IdP) even if it has a valid session to the IdP |
| `entityId`                  | If you have to use a customized Entity Id, you can overwrite it here. Normally, the default value can be left as it is |
| `loginFailedUrl`            | The URL to redirect to if the login failed. If this value is null the default page /noLogin will be shown |

The configuration is done via a corresponding JSON structure. As an example, the following JSON structure can be used to define an endpoint for app permissions:

```json
{
  "id": 1,
  "enable": true,
  "serviceProviderCertificate": "-----BEGIN CERTIFICATE-----***-----END CERTIFICATE-----",
  "serviceProviderPrivateKey": "-----BEGIN PRIVATE KEY-----***-----END PRIVATE KEY-----",
  "idpMetadataUrl": "https://login.thirdparty-idp.example/federationmetadata/metadata.xml",
  "idpMetadataXml": null,
  "ownUrl": "https://my.eliona-instance.example",
  "userToArchive": false,
  "allowInitializationByIdp": false,
  "signedRequest": true,
  "forceAuthn": false,
  "entityId": "{ownUrl}/saml/metadata",
  "loginFailedUrl": "https://myFancyLogoutPage.example"
}
```

Configurations can be created using this structure in Eliona under `Apps > Exchange app > Settings`. To do this, select the /configs endpoint with the POST method.

After completing configuration, the app starts Continuous Asset Creation. When all discovered rooms are created, user is notified about that in Eliona's notification system.

For detailed configuration steps, refer to your SSO provider's documentation.

### Microsoft ADFS Settings

To configure Microsoft ADFS specifically, follow these steps:

1. **Register a New App in Azure**
   - Go to the [Azure portal](https://portal.azure.com/).
   - Navigate to **Azure Active Directory** > **App registrations** > **New registration**.
   - Enter your application name and redirect URI (e.g., `https://customer.eliona.cloud/apps-public/saml-sso/`, platform `Web`).
   - Click **Register**.
   ![Azure app registration](user_guide/azure_app_registration.avif)

2. **Generate a Client Secret**
   - Go to **Certificates & secrets**.
   - Under **Client secrets**, click **New client secret**.
   - Add a description and set an expiration period, then click **Add**.
   - Copy the value of the client secret and store it securely. You will need this for Eliona configuration, and you won't be able to access it later.

3. **Locate Configuration Data**
   - Find the **Application (client) ID** and **Directory (tenant) ID** on the app's **Overview** page.

4. **Configure ADFS Settings in Eliona**
   - **MS Log-in**: Activate the log-in button "via Microsoft" by setting the configuration as "Enabled".
   - **Metadata URL**: Enter the Metadata URL from your Microsoft Azure account (found under app registration -> Endpoints).
     ![Metadata URL](user_guide/metadata.png)
   - **Own URL**: Enter your Eliona system URL (e.g., `https://customer.eliona.cloud`).
   - **Private Key**: Enter the private key in PEM format, matching your Azure certificate (found under Certificates & secrets -> Certificate).
     ![Private Key](user_guide/certificate.png)
   - **Certificate**: Can be a self-generated certificate.

For detailed steps on how to register an app in Azure, refer to the official [Microsoft documentation](https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-register-app).

For more information on generating and managing certificates, see the [Azure Key Vault documentation](https://docs.microsoft.com/en-us/azure/key-vault/certificates/).
