function sendPostRequestRedirect(req_data, req_url, redir_path) {
  $.ajax({
    url: req_url,
    dataType: 'json',
    type: 'POST',
    contentType: 'application/json',
    data: JSON.stringify(req_data),
    processData: false,
    success: function (data, textStatus, jQxhr) {
      console.log('OK');
      window.location.href = redir_path;
    },
    error: function (jqXhr, textStatus, errorThrown) {
      console.log('FUCK');
    }
  });
}

function AgentGenericAction(_action, _agents) {
  var _url = '/agents/' + _action;
  var data = {
    action: _action,
    agents: _agents
  };
  sendPostRequestRedirect(data, _url, "/");
}

function AgentDeleteAction(_agents) {
  AgentGenericAction('AGENT_DELETE', _agents)
}

function AgentHideAction(_agents) {
  AgentGenericAction('AGENT_HIDE', _agents)
}

function CallbackRemove(_callbackids) {
  var _url = '/callbacks/CALLBACK_DELETE';
  var data = {
    action: 'CALLBACK_DELETE',
    callbacks: _callbackids
  };
  sendPostRequestRedirect(data, _url, '/callbacks');
}
