package constants

// Dirs
const TemplatesDir = "templates"
const ConfigDir = "./config"
const ServiceDirPlaceholder = `_servicename_`

// Files
const TemplateExtension = ".tmpl"
const ServiceFilePlaceholder = `_service_`
const HandlerPlaceHolder = `_handler_`

// Messages
const ProjectBuiltMsg = `Built %s, don't forget to commit`
const UndefinedCommandMsg = `Cannot parse command %s`

// HTTP
const PayloadLimit = "limit"
const PayloadOffset = "offset"
const ContentJson = "application/json"
const APIHTTPAdapter = "api/http"

// Storage
const Storage = "storage"
const StorageEntityPlaceholder = "_entity_"
const StorageEntityOpPlaceholder = "_entity_op_"

// To-do(rolorin): move this to cli ops
const DefaultStorage = "mongo"

// misc
const Service = "service"
