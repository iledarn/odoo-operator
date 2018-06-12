package cluster

import (
	"bytes"
	"fmt"
)

const (
	// Volume Names
	configVolName = "config"

	// Ports and Port Names
	clientPortName      = "client-port"
	clientPort          = 8069
	longpollingPortName = "lp-port"
	longpollingPort     = 8072
)

const (
	odooConfigDir     = "/opt/odoo/odoorc.d"
	odooDefaultConfig = "01-default"
	odooCustomConfig  = "02-custom"
	// Basic Config
	odooPersistenceDir        = "/var/lib/odoo-persist"
	odooWithoutDemo           = "True"
	odooServerWideModules     = "web,web_kanban,backup_all"
	odooDbName                = "False"
	odooDbTemplate            = "template1"
	odooListDb                = "False"
	odooDbFilter              = "^%h$"
	odooBackupDir             = "/var/lib/odoo-backups"
	odooIntegratorWarrantyURL = "https://xoe.solutions/integrator-warranty/"
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

const odooBasicFmt = `
[options]

;; ========== BASIC CONFIG ==========
data_dir = %s

;; ========== ADVANCED CONFIG =======
without_demo = %s
; import_partial =
; pidfile =
server_wide_modules = %s

;; ========== HTTP/RPC SETTING ======
proxy_mode = True
; xmlrpc_interface =
; xmlrpc_port = 8069
; xmlrpc = True
; longpolling_port = 8072

;; ========== DATABASE ==============
db_name = %s
; pg_path =
db_template = %s

;; ========== SECURITY CONFIG =======
;; Through API
list_db = %s
;; Through Interface
dbfilter = %s

;; ========== ADVANCED CONFIG =======
; osv_memory_count_limit = False
; osv_memory_age_limit = 1.0
unaccent = True
; geoip_database = /usr/share/GeoIP/GeoLiteCity.dat
; csv_internal_sep = ,
; publisher_warranty_url = http://services.openerp.com/publisher-warranty/
; reportgz = False
; shell_interface = [ipython|ptpython|bpython|python]

;; ========== CUSTOM SECTIONS =======
[backups]
backupfolder = %s

[integrator]
integrator_warranty_url = %s

;; ==================================
;; ==================================
`

const odooLoggingFmt = `
[options]

;; ========== LOGGING SETTING =======
; logfile = False
; logrotate = False
; syslog = False
log_handler = %s
; log_db = False
; log_db_level = warning

;; ==================================
;; ==================================
`

const odooMultiprocFmt = `
[options]

;; ========== MULTI PROCESSING ======
max_cron_threads = 1
workers = 8
; limit_memory_soft = 2048 * 1024 * 1024
; limit_memory_hard = 2560 * 1024 * 1024
limit_time_cpu = 3600
limit_time_real = 240
limit_time_real_cron = 360
; limit_request = 8192

db_maxconn = 14

;; ==================================
;; ==================================
`

const odooSMTPFmt = `
[options]
;; ==================================
;; ========== SMTP SETTING ==========
email_from = %s
smtp_server = %s
smtp_port = %s
smtp_ssl = %s
smtp_user = %s
smtp_password = %s
`

func newConfigWithDefaultParams(data string) string {
	buf := bytes.NewBufferString(data)
	basicSection := fmt.Sprintf(odooBasicFmt,
		odooPersistenceDir,
		odooWithoutDemo,
		odooServerWideModules,
		odooDbName,
		odooDbTemplate,
		odooListDb,
		odooDbFilter,
		odooBackupDir,
		odooIntegratorWarrantyURL)
	buf.WriteString(basicSection)

	loggingSection := fmt.Sprintf(odooLoggingFmt,
		odooLogLevel)
	buf.WriteString(loggingSection)

	// multiprocSection := fmt.Sprintf(odooMultiprocFmt,
	// 	"")
	buf.WriteString(odooMultiprocFmt)

	SMTPSection := fmt.Sprintf(odooSMTPFmt,
		odooSMTPMail,
		odooSMTPServer,
		odooSMTPPort,
		odooSMTPSsl,
		odooSMTPUser,
		odooSMTPPassword)
	buf.WriteString(SMTPSection)

	return buf.String()
}
