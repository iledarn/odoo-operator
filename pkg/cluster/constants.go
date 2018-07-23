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
	envPGHOST         = "PGHOST"
	envPGPASSFILE     = "PGPASSFILE"
	envODOORC         = "ODOO_RC"
	envODOOPASSFILE   = "ODOO_PASSFILE"
	odooConfigDir     = "/opt/odoo/odoorc.d/"
	odooDefaultConfig = "01-default"
	odooCustomConfig  = "02-custom"
	odooSecretDir     = "/run/secrets/odoo/"
	odooPsqlSecret    = "pgpass"
	odooAdminSecret   = "adminpwd"
	// Basic Config
	odooVolumeMountPath       = "/mnt/odoo/"
	odooPersistenceDir        = odooVolumeMountPath + "persist/"
	odooWithoutDemo           = "True"
	odooServerWideModules     = "web,web_kanban,backup_all"
	odooDbName                = "False"
	odooDbTemplate            = "template1"
	odooListDb                = "False"
	odooDbFilter              = "^%h$"
	odooBackupDir             = odooVolumeMountPath + "backups/"
	odooIntegratorWarrantyURL = "https://erp.xoe.solutions/integrator-warranty/"
	// Log Config
	odooLogLevel = ":INFO"
	// Multiproc Config
	// SMTP Server Config
	odooSMTPMail     = ""
	odooSMTPServer   = ""
	odooSMTPPort     = ""
	odooSMTPSsl      = ""
	odooSMTPUser     = ""
	odooSMTPPassword = ""
)
