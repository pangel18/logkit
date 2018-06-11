package sender

import (
	. "github.com/qiniu/logkit/utils/models"
)

// ModeUsages 用途说明
var ModeUsages = []KeyValue{
	{TypePandora, "发送至 七牛云智能日志平台(Pandora)"},
	{TypeFile, "发送至 本地文件"},
	{TypeMongodbAccumulate, "发送至 MongoDB 服务"},
	{TypeInfluxdb, "发送至 InfluxDB 服务"},
	{TypeDiscard, "消费数据但不发送"},
	{TypeElastic, "发送至 Elasticsearch 服务"},
	{TypeKafka, "发送至 Kafka 服务"},
	{TypeHttp, "发送至 HTTP 服务器"},
}

var (
	OptionSaveLogPath = Option{
		KeyName:      KeyFtSaveLogPath,
		ChooseOnly:   false,
		Default:      "",
		DefaultNoUse: false,
		Description:  "管道磁盘数据保存路径(ft_save_log_path)",
		Advance:      true,
		ToolTip:      `指定备份数据的存放路径`,
	}
	OptionFtWriteLimit = Option{
		KeyName:      KeyFtWriteLimit,
		ChooseOnly:   false,
		Default:      "",
		DefaultNoUse: false,
		Description:  "磁盘写入限速(ft_write_limit)",
		CheckRegex:   "\\d+",
		Advance:      true,
		ToolTip:      `为了避免速率太快导致磁盘压力加大，可以根据系统情况自行限定写入本地磁盘的速率，单位MB/s`,
	}
	OptionFtStrategy = Option{
		KeyName:       KeyFtStrategy,
		ChooseOnly:    true,
		ChooseOptions: []interface{}{KeyFtStrategyBackupOnly, KeyFtStrategyAlwaysSave, KeyFtStrategyConcurrent},
		Default:       KeyFtStrategyBackupOnly,
		DefaultNoUse:  false,
		Description:   "磁盘管道容错策略[仅备份错误|全部数据走管道|仅增加并发](ft_strategy)",
		Advance:       true,
		ToolTip:       `设置为backup_only的时候，数据不经过本地队列直接发送到下游，设为always_save时则所有数据会先发送到本地队列，选concurrent的时候会直接并发发送，不经过队列。无论该选项设置什么，失败的数据都会加入到重试队列中异步循环重试`,
	}
	OptionFtProcs = Option{
		KeyName:      KeyFtProcs,
		ChooseOnly:   false,
		Default:      "",
		DefaultNoUse: false,
		Description:  "发送并发数量(ft_procs)",
		CheckRegex:   "\\d+",
		Advance:      true,
		ToolTip:      "并发仅在ft_strateg模式选择 always_save或concurrent 时生效",
	}
	OptionFtMemoryChannel = Option{
		KeyName:       KeyFtMemoryChannel,
		Element:       Radio,
		ChooseOnly:    true,
		ChooseOptions: []interface{}{"false", "true"},
		Default:       "false",
		DefaultNoUse:  false,
		Description:   "用内存管道(ft_memory_channel)",
		Advance:       true,
		ToolTip:       `内存管道替代磁盘管道`,
	}
	OptionFtMemoryChannelSize = Option{
		KeyName:       KeyFtMemoryChannelSize,
		ChooseOnly:    false,
		Default:       "",
		DefaultNoUse:  false,
		Description:   "内存管道长度(ft_memory_channel_size)",
		CheckRegex:    "\\d+",
		Advance:       true,
		AdvanceDepend: KeyFtMemoryChannel,
		ToolTip:       `默认为"100"，单位为批次，也就是100代表100个待发送的批次，注意：该选项设置的大小表达的是队列中可存储的元素个数，并不是占用的内存大小`,
	}
	OptionLogkitSendTime = Option{
		KeyName:       KeyLogkitSendTime,
		Element:       Radio,
		ChooseOnly:    true,
		ChooseOptions: []interface{}{"true", "false"},
		Default:       "true",
		DefaultNoUse:  false,
		Description:   "添加系统时间(logkit_send_time)",
		Advance:       true,
		ToolTip:       "在系统中添加数据发送时的当前时间作为时间戳",
	}
)
var ModeKeyOptions = map[string][]Option{
	TypeFile: {
		{
			KeyName:      KeyFileSenderPath,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "/home/john/mylogs/my-%Y-%m-%d.log",
			DefaultNoUse: true,
			Description:  "发送到指定文件(file_send_path)",
			ToolTip:      `路径支持魔法变量，例如 "file_send_path":"data-%Y-%m-%d.txt" ，此时数据就会渲染出日期，存放为 data-2018-03-28.txt`,
		},
	},
	TypePandora: {
		{
			KeyName:      KeyPandoraWorkflowName,
			ChooseOnly:   false,
			Default:      "",
			Placeholder:  "logkit_default_workflow",
			DefaultNoUse: true,
			Required:     true,
			Description:  "工作流名称(pandora_workflow_name)",
			CheckRegex:   "^[a-zA-Z_][a-zA-Z0-9_]{0,127}$",
			ToolTip:      "七牛大数据平台工作流名称",
		},
		{
			KeyName:      KeyPandoraRepoName,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "my_work",
			DefaultNoUse: true,
			Description:  "数据源名称(pandora_repo_name)",
			CheckRegex:   "^[a-zA-Z][a-zA-Z0-9_]{0,127}$",
			ToolTip:      "七牛大数据平台工作流中的数据源名称",
		},
		{
			KeyName:      KeyPandoraAk,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "在此填写您七牛账号ak(access_key)",
			DefaultNoUse: true,
			Description:  "七牛公钥(access_key)",
		},
		{
			KeyName:      KeyPandoraSk,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "在此填写您七牛账号的sk(secret_key)",
			DefaultNoUse: true,
			Description:  "七牛私钥(secret_key)",
			Secret:       true,
		},
		OptionSaveLogPath,
		OptionLogkitSendTime,
		{
			KeyName:      KeyPandoraHost,
			ChooseOnly:   false,
			Default:      "https://pipeline.qiniu.com",
			DefaultNoUse: false,
			Description:  "大数据平台域名(pandora_host)",
			Advance:      true,
			ToolTip:      "数据发送的目的域名，私有部署请对应修改",
		},
		{
			KeyName:       KeyPandoraRegion,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"nb"},
			Default:       "nb",
			DefaultNoUse:  false,
			Description:   "创建的资源所在区域(pandora_region)",
			Advance:       true,
			ToolTip:       "工作流资源创建所在区域",
		},
		{
			KeyName:       KeyPandoraSchemaFree,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "自动创建数据源并更新(pandora_schema_free)",
			Advance:       true,
			ToolTip:       "自动根据数据创建工作流、数据源并自动更新",
		},
		{
			KeyName:       KeyPandoraAutoCreate,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "以DSL语法自动创建数据源(pandora_auto_create)",
			Advance:       true,
			ToolTip:       `自动创建数据源，语法为 "f1 date, f2 string, f3 float, f4 map{f5 long}"`,
			ToolTipActive: true,
		},
		{
			KeyName:      KeyPandoraSchema,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "筛选字段(重命名)发送(pandora_schema)",
			Advance:      true,
			ToolTip:      `将f1重命名为f2: "f1 f2,...", 其他自动不要去掉...`,
		},
		{
			KeyName:       KeyPandoraEnableLogDB,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "自动创建并导出到日志分析(pandora_enable_logdb)",
		},
		{
			KeyName:       KeyPandoraLogDBName,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定日志分析仓库名称(pandora_logdb_name)",
			AdvanceDepend: KeyPandoraEnableLogDB,
			ToolTip:       "若不指定使用数据源(pandora_repo_name)名称",
		},
		{
			KeyName:       KeyPandoraLogDBHost,
			ChooseOnly:    false,
			Default:       "https://logdb.qiniu.com",
			DefaultNoUse:  false,
			Description:   "日志分析域名[私有部署才修改](pandora_logdb_host)",
			Advance:       true,
			AdvanceDepend: KeyPandoraEnableLogDB,
			ToolTip:       "日志分析仓库域名，私有部署请对应修改",
		},
		{
			KeyName:       KeyPandoraLogDBAnalyzer,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定字段分词方式(pandora_logdb_analyzer)",
			Advance:       true,
			AdvanceDepend: KeyPandoraEnableLogDB,
			ToolTip:       `指定字段的分词方式，逗号分隔多个，如 "f1 keyword, f2 full_text"`,
		},
		{
			KeyName:       KeyPandoraEnableTSDB,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "自动创建并导出到时序数据库(pandora_enable_tsdb)",
		},
		{
			KeyName:       KeyPandoraTSDBName,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定时序数据库仓库名称(pandora_tsdb_name)",
			AdvanceDepend: KeyPandoraEnableTSDB,
			ToolTip:       "若不指定使用数据源(pandora_repo_name)名称",
		},
		{
			KeyName:       KeyPandoraTSDBSeriesName,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定时序数据库序列名称(pandora_tsdb_series_name)",
			AdvanceDepend: KeyPandoraEnableTSDB,
			ToolTip:       "若不指定使用仓库(pandora_tsdb_name)名称",
		},
		{
			KeyName:       KeyPandoraTSDBSeriesTags,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定时序数据库标签(pandora_tsdb_series_tags)",
			AdvanceDepend: KeyPandoraEnableTSDB,
		},
		{
			KeyName:       KeyPandoraTSDBHost,
			ChooseOnly:    false,
			Default:       "https://tsdb.qiniu.com",
			DefaultNoUse:  false,
			Description:   "时序数据库域名[私有部署才修改](pandora_tsdb_host)",
			Advance:       true,
			AdvanceDepend: KeyPandoraEnableTSDB,
			ToolTip:       "时序数据库域名，私有部署请对应修改",
		},
		{
			KeyName:       KeyPandoraTSDBTimeStamp,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "指定时序数据库时间戳(pandora_tsdb_timestamp)",
			Advance:       true,
			AdvanceDepend: KeyPandoraEnableTSDB,
		},
		{
			KeyName:       KeyPandoraEnableKodo,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "自动导出到七牛云存储(pandora_enable_kodo)",
		},
		{
			KeyName:       KeyPandoraKodoBucketName,
			ChooseOnly:    false,
			Default:       "",
			Required:      true,
			Placeholder:   "my_bucket_name",
			DefaultNoUse:  true,
			Description:   "云存储仓库名称(启用自动导出到云存储时必填)(pandora_bucket_name)",
			AdvanceDepend: KeyPandoraEnableKodo,
		},
		{
			KeyName:       KeyPandoraEmail,
			ChooseOnly:    false,
			Default:       "",
			Required:      true,
			Placeholder:   "my@email.com",
			DefaultNoUse:  true,
			Description:   "七牛账户邮箱(qiniu_email)",
			AdvanceDepend: KeyPandoraEnableKodo,
		},
		{
			KeyName:       KeyPandoraKodoFilePrefix,
			ChooseOnly:    false,
			Default:       "logkitauto/date=$(year)-$(mon)-$(day)/hour=$(hour)/min=$(min)/$(sec)",
			DefaultNoUse:  false,
			Description:   "云存储文件前缀(pandora_kodo_prefix)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraKodoCompressPrefix,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"parquet", "json", "text", "csv"},
			Default:       "parquet",
			DefaultNoUse:  false,
			Description:   "云存储文件保存格式(pandora_kodo_compress)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraKodoGzip,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "云存储文件压缩存储(pandora_kodo_gzip)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraKodoRotateStrategy,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"interval", "size", "both"},
			Default:       "interval",
			DefaultNoUse:  false,
			Description:   "云存储文件分割策略(pandora_kodo_rotate_strategy)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraKodoRotateInterval,
			ChooseOnly:    false,
			Default:       "600",
			DefaultNoUse:  false,
			Description:   "云存储文件分割间隔时间(pandora_kodo_rotate_interval)(单位秒)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraKodoRotateSize,
			ChooseOnly:    false,
			Default:       "512000",
			DefaultNoUse:  false,
			Description:   "云存储文件分割文件大小(pandora_kodo_rotate_size)(单位KB)",
			AdvanceDepend: KeyPandoraEnableKodo,
			Advance:       true,
		},
		{
			KeyName:       KeyPandoraGzip,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "压缩发送(pandora_gzip)",
			Advance:       true,
			ToolTip:       "使用gzip压缩发送",
		},
		{
			KeyName:      KeyFlowRateLimit,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "流量限制(flow_rate_limit)",
			CheckRegex:   "\\d+",
			Advance:      true,
			ToolTip:      "对请求流量限制,单位为[KB/s]",
		},
		{
			KeyName:      KeyRequestRateLimit,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "请求限制(request_rate_limit)",
			CheckRegex:   "\\d+",
			Advance:      true,
			ToolTip:      "对请求次数限制,单位为[次/s]",
		},
		{
			KeyName:       KeyPandoraUUID,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "数据植入UUID(pandora_uuid)",
			Advance:       true,
			ToolTip:       `该字段保证了发送出去的每一条数据都拥有一个唯一的UUID，可以用于数据去重等需要`,
		},
		{
			KeyName:       KeyPandoraWithIP,
			ChooseOnly:    false,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "数据植入来源IP(pandora_withip)",
			Advance:       true,
		},
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
		{
			KeyName:       KeyForceMicrosecond,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "扰动时间字段增加精度(force_microsecond)",
			Advance:       true,
		},
		{
			KeyName:       KeyForceDataConvert,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "自动转换类型(pandora_force_convert)",
			Advance:       true,
			ToolTip:       `强制类型转换，如定义的pandora schema为long，而实际的为string，则会尝试将string解析为long类型，若无法解析，则忽略该字段内容`,
		},
		{
			KeyName:       KeyNumberUseFloat,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "数字统一为float(number_use_float)",
			Advance:       true,
			ToolTip:       `对于整型和浮点型都自动识别为浮点型`,
		},
		{
			KeyName:       KeyIgnoreInvalidField,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "忽略错误字段(ignore_invalid_field)",
			Advance:       true,
			ToolTip:       `进行数据格式校验，并忽略不符合格式的字段数据`,
		},
		{
			KeyName:       KeyPandoraAutoConvertDate,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "自动转换时间类型(pandora_auto_convert_date)",
			Advance:       true,
			ToolTip:       `会自动将用户的自动尝试转换为Pandora的时间类型(date)`,
		},
		{
			KeyName:       KeyPandoraUnescape,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  false,
			Description:   "服务端反转译换行/制表符(pandora_unescape)",
			Advance:       true,
			ToolTip:       `在pandora服务端反转译\\n=>\n, \\t=>\t; 由于pandora数据上传的编码方式会占用\t和\n这两个符号，所以在sdk中打点时会默认把\t转译为\\t，把\n转译为\\n，开启这个选项就是在服务端把这个反转译回来。开启该选项也会转译数据中原有的这些\\t和\\n符号`,
		},
		{
			KeyName:       KeyInsecureServer,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"false", "true"},
			Default:       "false",
			DefaultNoUse:  false,
			Description:   "认证不验证(insecure_server)",
			Advance:       true,
			ToolTip:       `对于https等情况不对证书和安全性检验`,
		},
	},
	TypeMongodbAccumulate: {
		{
			KeyName:      KeyMongodbHost,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]",
			DefaultNoUse: true,
			Description:  "数据库地址(mongodb_host)",
			ToolTip:      `Mongodb的地址: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]`,
		},
		{
			KeyName:      KeyMongodbDB,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "app123",
			DefaultNoUse: true,
			Description:  "数据库名称(mongodb_db)",
		},
		{
			KeyName:      KeyMongodbCollection,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "collection1",
			DefaultNoUse: true,
			Description:  "数据表名称(mongodb_collection)",
		},
		{
			KeyName:      KeyMongodbUpdateKey,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "domain,uid",
			DefaultNoUse: true,
			Description:  "聚合条件列(mongodb_acc_updkey)",
		},
		{
			KeyName:      KeyMongodbAccKey,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "low,hit",
			DefaultNoUse: true,
			Description:  "聚合列(mongodb_acc_acckey)",
		},
		OptionSaveLogPath,
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
	},
	TypeInfluxdb: {
		{
			KeyName:      KeyInfluxdbHost,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "127.0.0.1:8086",
			DefaultNoUse: true,
			Description:  "数据库地址(influxdb_host)",
			ToolTip:      `数据库地址127.0.0.1:8086`,
		},
		{
			KeyName:      KeyInfluxdbDB,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "testdb",
			DefaultNoUse: true,
			Description:  "数据库名称(influxdb_db)",
		},
		{
			KeyName:       KeyInfluxdbAutoCreate,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			Description:   "自动创建数据库(influxdb_auto_create)",
			Advance:       true,
		},
		{
			KeyName:      KeyInfluxdbMeasurement,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "test_table",
			DefaultNoUse: true,
			Description:  "measurement名称(influxdb_measurement)",
		},
		{
			KeyName:      KeyInfluxdbRetetion,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "retention名称(influxdb_retention)",
			Advance:      true,
		},
		{
			KeyName:       KeyInfluxdbRetetionDuration,
			ChooseOnly:    false,
			Default:       "",
			DefaultNoUse:  false,
			Description:   "retention时长(influxdb_retention_duration)",
			AdvanceDepend: KeyInfluxdbAutoCreate,
			Advance:       true,
		},
		{
			KeyName:      KeyInfluxdbTags,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "标签列数据(influxdb_tags)",
			Advance:      true,
		},
		{
			KeyName:      KeyInfluxdbFields,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "普通列数据(influxdb_fields)",
			Advance:      true,
		},
		{
			KeyName:      KeyInfluxdbTimestamp,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "时间戳列(influxdb_timestamp)",
			Advance:      true,
		},
		{
			KeyName:      KeyInfluxdbTimestampPrecision,
			ChooseOnly:   false,
			Default:      "100",
			DefaultNoUse: false,
			Description:  "时间戳列精度调整(influxdb_timestamp_precision)",
			Advance:      true,
		},
		OptionSaveLogPath,
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
	},
	TypeDiscard: {},
	TypeElastic: {
		{
			KeyName:      KeyElasticHost,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "192.168.31.203:9200",
			DefaultNoUse: false,
			Description:  "host地址(elastic_host)",
			ToolTip:      `常用端口9200`,
		},
		{
			KeyName:       KeyElasticVersion,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{ElasticVersion3, ElasticVersion5, ElasticVersion6},
			Description:   "ES版本号(es_version)",
		},
		{
			KeyName:      KeyElasticIndex,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "app-repo-123",
			DefaultNoUse: true,
			Description:  "索引名称(elastic_index)",
		},
		{
			KeyName:       KeyElasticIndexStrategy,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{KeyDefaultIndexStrategy, KeyYearIndexStrategy, KeyMonthIndexStrategy, KeyDayIndexStrategy},
			Default:       KeyFtStrategyBackupOnly,
			DefaultNoUse:  false,
			Description:   "自动索引模式(默认索引|按年索引|按月索引|按日索引)(index_strategy)",
			Advance:       true,
		},
		{
			KeyName:       KeyElasticTimezone,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{KeyUTCTimezone, KeylocalTimezone, KeyPRCTimezone},
			Default:       KeyUTCTimezone,
			DefaultNoUse:  false,
			Description:   "索引时区(Local(本地)|UTC(标准时间)|PRC(北京时间))(elastic_time_zone)",
			Advance:       true,
		},
		OptionLogkitSendTime,
		{
			KeyName:      KeyElasticType,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "app",
			DefaultNoUse: true,
			Description:  "索引类型名称(elastic_type)",
		},
		OptionSaveLogPath,
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
	},
	TypeKafka: {
		{
			KeyName:      KeyKafkaHost,
			ChooseOnly:   false,
			Required:     true,
			Default:      "",
			Placeholder:  "192.168.31.201:9092",
			DefaultNoUse: true,
			Description:  "broker的host地址(kafka_host)",
			ToolTip:      "常用端口 9092",
		},
		{
			KeyName:      KeyKafkaTopic,
			ChooseOnly:   false,
			Default:      "",
			Required:     true,
			Placeholder:  "my_topic",
			DefaultNoUse: true,
			Description:  "打点的topic名称(kafka_topic)",
		},
		{
			KeyName:       KeyKafkaCompression,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{KeyKafkaCompressionNone, KeyKafkaCompressionGzip, KeyKafkaCompressionSnappy},
			Default:       KeyKafkaCompressionNone,
			DefaultNoUse:  false,
			Description:   "压缩模式[none不压缩|gzip压缩|snappy压缩](kafka_compression)",
		},
		{
			KeyName:      KeyKafkaClientId,
			ChooseOnly:   false,
			Default:      "",
			DefaultNoUse: false,
			Description:  "kafka客户端标识ID(kafka_client_id)",
			Advance:      true,
		},
		{
			KeyName:      KeyKafkaRetryMax,
			ChooseOnly:   false,
			Default:      "3",
			DefaultNoUse: false,
			Description:  "kafka最大错误重试次数(kafka_retry_max)",
			Advance:      true,
		},
		{
			KeyName:      KeyKafkaTimeout,
			ChooseOnly:   false,
			Default:      "30s",
			DefaultNoUse: false,
			Description:  "kafka连接超时时间(kafka_timeout)",
			Advance:      true,
		},
		{
			KeyName:      KeyKafkaKeepAlive,
			ChooseOnly:   false,
			Default:      "0",
			DefaultNoUse: false,
			Description:  "kafka的keepalive时间(kafka_keep_alive)",
			Advance:      true,
		},
		OptionSaveLogPath,
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
	},
	TypeHttp: {
		{
			KeyName:      KeyHttpSenderUrl,
			ChooseOnly:   false,
			Default:      "",
			Placeholder:  "http://127.0.0.1/data",
			DefaultNoUse: true,
			Required:     true,
			Description:  "发送目的url(http_sender_url)",
		},
		{
			KeyName:       KeyHttpSenderProtocol,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"json", "csv"},
			Default:       "json",
			Description:   "发送数据时使用的格式(http_sender_protocol)",
		},
		{
			KeyName:      KeyHttpSenderCsvSplit,
			ChooseOnly:   false,
			Default:      "",
			Placeholder:  ",",
			Required:     true,
			DefaultNoUse: true,
			Description:  "csv分隔符(http_sender_csv_split)",
		},
		{
			KeyName:       KeyHttpSenderGzip,
			Element:       Radio,
			ChooseOnly:    true,
			ChooseOptions: []interface{}{"true", "false"},
			Default:       "true",
			DefaultNoUse:  true,
			Description:   "是否启用gzip(http_sender_gzip)",
		},
		OptionSaveLogPath,
		OptionFtWriteLimit,
		OptionFtStrategy,
		OptionFtProcs,
		OptionFtMemoryChannel,
		OptionFtMemoryChannelSize,
	},
}
