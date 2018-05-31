package cluster

import (
	"bytes"
	"fmt"
	"path/filepath"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

const (
	odooDefaultConfigPath = "/opt/odoo/odoorc.d/01-default"
	odooCustomConfigPath  = "/opt/odoo/odoorc.d/02-custom"
	// Basic Config
	odooDataDir               = "/var/lib/odoo-persist"
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

func configmapForOdooCluster(cr *api.OdooCluster) *v1.ConfigMap {

	var cfgDefaultData string
	var cfgCustomData string

	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
		},
	}
	cm.Name = "abc"
	cm.Labels = labelsForOdooCluster(cr.Name)
	cfgDefaultData = newConfigWithDefaultParams(cfgDefaultData)
	cm.Data = map[string]string{filepath.Base(odooDefaultConfigPath): cfgDefaultData}
	if len(cr.Spec.ConfigMap) != 0 {
		cfgCustomData = cr.Spec.ConfigMap
		cm.Data[filepath.Base(odooCustomConfigPath)] = cfgCustomData
	}
	addOwnerRefToObject(cm, asOwner(cr))

	return cm
}

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
		odooDataDir,
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
