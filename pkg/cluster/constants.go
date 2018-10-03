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
	appDefaultConfigKey    = "01-DEFAULT"
	appOptionsConfigKey    = "02-options"
	appIntegratorConfigKey = "03-integrator"
	appCustomConfigKey     = "04-custom"
	appPsqlSecretKey       = "pgpass"
	appAdminSecretKey      = "adminpwd"

	// Basic Config
	defaultServerTierMaxConn      = "16"
	defaultLongpollingTierMaxConn = "16"
	defaultWithoutDemo            = "True"
	defaultServerWideModules      = "base,web"
	defaultDbName                 = "False"
	defaultDbTemplate             = "template0"
	defaultListDb                 = "False"
	defaultDbFilter               = "^%h$"
	defaultPublisherWarrantyURL   = "http://services.openerp.com/publisher-warranty/"

	// Log Config
	defaultLogLevel = ":INFO"
)
