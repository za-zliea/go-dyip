export default {
  brand: {
    name: 'ZlieA'
  },
  console: {
    title: 'ZlieA DDNS Management'
  },
  common: {
    cancel: 'Cancel',
    confirm: 'Confirm',
    refresh: 'Refresh',
    operation: 'Action',
    detail: 'Detail',
    submit: 'Submit',
    empty: 'No data',
    close: 'Close'
  },
  route: {
    login: 'Login',
    ddnsView: 'DDNS View',
    ddnsUpdate: 'DDNS Update'
  },
  lang: {
    label: 'Language',
    zh: '中文',
    en: 'English'
  },
  login: {
    username: 'Username',
    password: 'Password',
    placeholderUsername: 'Enter your username',
    placeholderPassword: 'Enter your password',
    button: 'Sign In',
    success: 'Signed in successfully',
    failed: 'Sign-in failed: incorrect username or password',
    ruleUsername: 'Please enter your username',
    rulePassword: 'Please enter your password'
  },
  layout: {
    logout: 'Sign Out',
    logoutButton: 'Sign Out',
    logoutConfirmTitle: 'Notice',
    logoutConfirmPrompt: 'Are you sure you want to sign out?'
  },
  ddnsView: {
    title: 'DDNS View',
    refresh: 'Refresh',
    empty: 'No domain records',
    loadingDns: 'Resolving',
    columnName: 'Domain',
    columnProvider: 'Provider',
    columnProtocol: 'Protocol',
    columnConsoleEnabled: 'Console update',
    columnClientUpload: 'Client upload',
    columnIp: 'Latest IP',
    columnUpdateTime: 'Updated',
    columnDip: 'DNS IP',
    columnConsistent: 'Consistent',
    columnOperation: 'Action',
    tagProtocolV4: 'IPv4',
    tagProtocolV6: 'IPv6',
    tagConsoleEnabled: 'Supported',
    tagConsoleDisabled: 'Not supported',
    tagClientUploadYes: 'Yes',
    tagClientUploadNo: 'No',
    tagConsistentYes: 'Consistent',
    tagConsistentNo: 'Inconsistent',
    tagConsistentUnknown: 'Unknown',
    btnDetail: 'Detail',
    btnUpdate: 'Update',
    detail: {
      title: 'Record detail',
      name: 'Domain',
      provider: 'Provider',
      protocol: 'Protocol',
      ip: 'Latest IP',
      updateTime: 'Updated',
      dip: 'DNS IP',
      consistent: 'Consistent',
      consoleEnabled: 'Console update',
      clientUpload: 'Client upload',
      historyTitle: 'Update history',
      columnHistoryIp: 'IP',
      columnHistoryTime: 'Time',
      emptyHistory: 'No history'
    }
  },
  ddnsUpdate: {
    title: 'DDNS Update',
    record: 'Record',
    placeholderRecord: 'Select a record',
    ruleRecord: 'Please select a record',
    ip: 'IP address',
    placeholderIp: 'Enter the IP address',
    ruleIp: 'Please enter the IP address',
    ruleIpFormat: 'Invalid IP address',
    currentIp: 'Current IP',
    useCurrentIp: 'Use current IP',
    currentIpMissing: 'Could not detect your current IP',
    currentIpFamilyMismatch: 'Your current IP does not match the selected record protocol ({protocol})',
    submit: 'Submit',
    success: 'IP updated successfully',
    sameIp: 'IP unchanged, skipped',
    alertTitle: 'No records support console update',
    alertDesc: 'Only records with console update enabled (synctype first digit 1) can be updated manually here.'
  },
  request: {
    failed: 'Request failed',
    failedWithStatus: 'Request failed ({status})',
    expired: 'Session expired, please sign in again',
    network: 'Network error, please check the connection to the server',
    unknown: 'Unknown error'
  }
}
