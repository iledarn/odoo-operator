package cluster

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
publisher_warranty_url = %s
; reportgz = False
; shell_interface = [ipython|ptpython|bpython|python]
`

const odooLoggingFmt = `

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
; We don't use Odoo's multiprocessing in k8s deployments
; This is handed off to the microservice infrastructure
;; ========== MULTI PROCESSING ======
; max_cron_threads = 2
; workers = 0
; limit_memory_soft = 2048 * 1024 * 1024
; limit_memory_hard = 2560 * 1024 * 1024
; limit_time_cpu = 3600
; limit_time_real = 240
; limit_time_real_cron = 360
; limit_request = 8192

; db_maxconn = 64

;; ==================================
;; ==================================
`

const odooSMTPFmt = `
;; ==================================
;; ========== SMTP SETTING ==========
email_from = %s
smtp_server = %s
smtp_port = %s
smtp_ssl = %s
smtp_user = %s
smtp_password = %s
`

const odooCustomSection = `

;; ========== CUSTOM SECTIONS =======
[backups]
backupfolder = %s

[integrator]
integrator_warranty_url = %s

;; ==================================
;; ==================================
`

const odooPsqlSecretFmt = `
# hostname:port:database:username:password
# If an entry needs to contain : or \, escape this character with \
%s:%s:*:%s:%s
`
const odooAdminSecretFmt = `%s

`
