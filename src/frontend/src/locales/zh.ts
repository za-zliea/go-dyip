export default {
  brand: {
    name: 'ZlieA'
  },
  console: {
    title: 'ZlieA 动态DNS解析管理控制台'
  },
  common: {
    cancel: '取消',
    confirm: '确定',
    refresh: '刷新',
    operation: '操作',
    detail: '详情',
    submit: '提交',
    empty: '暂无数据',
    close: '关闭'
  },
  route: {
    login: '登录',
    ddnsView: 'DDNS查看',
    ddnsUpdate: 'DDNS更新'
  },
  lang: {
    label: '语言',
    zh: '中文',
    en: 'English'
  },
  login: {
    username: '用户名',
    password: '密码',
    placeholderUsername: '请输入用户名',
    placeholderPassword: '请输入密码',
    button: '登 录',
    success: '登录成功',
    failed: '登录失败，用户名或者密码错误',
    ruleUsername: '请输入用户名',
    rulePassword: '请输入密码'
  },
  layout: {
    logout: '退出登录',
    logoutButton: '退出',
    logoutConfirmTitle: '提示',
    logoutConfirmPrompt: '确定要退出登录吗？'
  },
  ddnsView: {
    title: 'DDNS查看',
    refresh: '刷新',
    empty: '暂无域名记录',
    loadingDns: '解析中',
    columnName: '完整域名',
    columnProvider: '提供商',
    columnProtocol: '协议',
    columnConsoleEnabled: '控制台更新',
    columnClientUpload: '客户端上传',
    columnIp: '记录最新IP',
    columnUpdateTime: '更新时间',
    columnDip: 'DNS对应IP',
    columnConsistent: '是否一致',
    columnOperation: '操作',
    tagProtocolV4: 'IPv4',
    tagProtocolV6: 'IPv6',
    tagConsoleEnabled: '支持',
    tagConsoleDisabled: '不支持',
    tagClientUploadYes: '是',
    tagClientUploadNo: '否',
    tagConsistentYes: '一致',
    tagConsistentNo: '不一致',
    tagConsistentUnknown: '未知',
    btnDetail: '详情',
    btnUpdate: '更新',
    detail: {
      title: '记录详情',
      name: '完整域名',
      provider: '提供商',
      protocol: '协议',
      ip: '记录最新IP',
      updateTime: '更新时间',
      dip: 'DNS对应IP',
      consistent: '是否一致',
      consoleEnabled: '控制台更新',
      clientUpload: '客户端上传',
      historyTitle: '更新历史',
      columnHistoryIp: 'IP',
      columnHistoryTime: '时间',
      emptyHistory: '暂无历史记录'
    }
  },
  ddnsUpdate: {
    title: 'DDNS更新',
    record: '选择记录',
    placeholderRecord: '请选择记录',
    ruleRecord: '请选择记录',
    ip: 'IP 地址',
    placeholderIp: '请输入 IP 地址',
    ruleIp: '请输入 IP 地址',
    ruleIpFormat: 'IP 地址格式不正确',
    currentIp: '当前IP',
    useCurrentIp: '使用当前IP',
    currentIpMissing: '未获取到当前IP',
    currentIpFamilyMismatch: '当前IP 与所选记录协议（{protocol}）不匹配',
    submit: '提交更新',
    success: 'IP 更新成功',
    sameIp: 'IP 未变化，已跳过',
    alertTitle: '当前没有可控制台更新的记录',
    alertDesc: '只有开启控制台更新（synctype 首位为 1）的记录支持在此手动更新 IP。'
  },
  request: {
    failed: '请求失败',
    failedWithStatus: '请求失败 ({status})',
    expired: '登录已过期，请重新登录',
    network: '网络异常，请检查与服务端的连接',
    unknown: '未知错误'
  }
}
