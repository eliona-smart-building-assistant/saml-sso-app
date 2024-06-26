# SAML SSO

## Overview

The SAML SSO (Security Assertion Markup Language Single Sign-On) app allows users to log into Eliona using various SSO providers, including Microsoft ADFS. This setup streamlines authentication by using a single set of credentials.

## Configuration

![Configuration frontend](user_guide/frontend.webp)

To integrate a generic SAML SSO provider with Eliona, follow these general steps:

### General SAML SSO Settings

1. **Enable SAML SSO**: Activate the log-in button "via SAML".
2. **Metadata URL**: Enter the Metadata URL provided by your SAML SSO provider.
3. **Own URL**: Enter your Eliona system URL (e.g., `https://customer.eliona.cloud`).
4. **Private Key**: Enter the private key in PEM format.
5. **Certificate**: Enter the certificate, which can be self-generated.

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
   - **MS Log-in**: Activate the log-in button "via Microsoft" by clicking "Active".
     ![MS Log-in](user_guide/login.avif)
   - **Metadata URL**: Enter the Metadata URL from your Microsoft Azure account (found under app registration -> Endpoints).
     ![Metadata URL](user_guide/metadata.png)
   - **Own URL**: Enter your Eliona system URL (e.g., `https://customer.eliona.cloud`).
   - **Private Key**: Enter the private key in PEM format, matching your Azure certificate (found under Certificates & secrets -> Certificate).
     ![Private Key](user_guide/certificate.png)
   - **Certificate**: Can be a self-generated certificate.

For detailed steps on how to register an app in Azure, refer to the official [Microsoft documentation](https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-register-app).

For more information on generating and managing certificates, see the [Azure Key Vault documentation](https://docs.microsoft.com/en-us/azure/key-vault/certificates/).
