package cluster

const (
	// Volume Names
	configVolName = "config"
	secretVolName = "secret"

	// Ports and Port Names
	clientPortName      string = "client-port"
	clientPort                 = 8069
	longpollingPortName        = "lp-port"
	longpollingPort            = 8072
)

const (
	// Environment Variables
	envPGHOST       = "PGHOST"
	envPGUSER       = "PGUSER"
	envPGPASSFILE   = "PGPASSFILE"
	envODOORC       = "ODOO_RC"
	envODOOPASSFILE = "ODOO_PASSFILE"

	// App paths
	appMountPath   = "/mnt/odoo/"
	appBasePath    = "/opt/odoo/"
	appSecretsPath = "/run/secrets/odoo/"
	appConfigsPath = "/run/configs/odoo/"

	// ConfigMaps, Secrets & Volumes Keys
	appDefaultConfigKey = "default"
	appCustomConfigKey  = "custom"
	appPsqlSecretKey    = "pgpass"
	appAdminSecretKey   = "adminpwd"

	// Basic Config
	odooServerTierMaxConn      = "16"
	odooLongpollingTierMaxConn = "16"
	odooWithoutDemo            = "True"
	odooServerWideModules      = "base,web"
	odooDbName                 = "False"
	odooDbTemplate             = "template0"
	odooListDb                 = "False"
	odooDbFilter               = "^%h$"
	odooPublisherWarrantyURL   = "http://services.openerp.com/publisher-warranty/"
	odooIntegratorWarrantyURL  = "https://erp.xoe.solutions/integrator-warranty/"

	// Log Config
	odooLogLevel = ":INFO"

	// Multiproc Config
	// SMTP Server Config
	odooSMTPMail     = ""
	odooSMTPServer   = ""
	odooSMTPPort     = "465"
	odooSMTPSsl      = "true"
	odooSMTPUser     = ""
	odooSMTPPassword = ""
)
