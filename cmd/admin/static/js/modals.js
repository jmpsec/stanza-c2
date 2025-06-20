function command_modal_code() {
  $("#commands").val('');
  $('#command_modal').modal('show');
  $('#command_modal .command_button').click(function () {
    var to_execute = $("#commands").val();
    var uuid = $("#agent_uuid").val();
    if (to_execute.length > 0) {
      C2executeCommand(uuid, to_execute);
    }
    $('#command_modal').modal('hide');
  });
  $('#command_modal').on('shown.bs.modal', function () {
    $('#commands').focus();
  });
  $(document).keypress(function (e) {
    if ($("#command_modal").hasClass('in') && (e.keycode == 13 || e.which == 13)) {
      $('#command_modal .command_button').click();
    }
  });
}

function get_modal_code() {
  $("#get_file").val('');
  $('#get_modal').modal('show');
  $('#get_modal .get_button').click(function () {
    var get_file = $("#get_file").val();
    var uuid = $("#agent_uuid").val();
    C2getFile(uuid, get_file);
    $('#get_modal').modal('hide');
  });
  $('#get_modal').on('shown.bs.modal', function () {
    $('#get_file').focus();
  });
  $(document).keypress(function (e) {
    if ($("#get_modal").hasClass('in') && (e.keycode == 13 || e.which == 13)) {
      $('#get_modal .get_button').click();
    }
  });
}

function put_modal_code() {
  $("#put_url").val('');
  $("#put_file").val('');
  $('#put_modal').modal('show');
  $('#put_modal .put_button').click(function () {
    var put_url = $("#put_url").val();
    var put_file = $("#put_file").val();
    var uuid = $("#agent_uuid").val();
    if (put_url.length > 0 && put_file.length > 0) {
      C2putFile(uuid, put_url, put_file);
    }
    $('#put_modal').modal('hide');
  });
  $('#put_modal').on('shown.bs.modal', function () {
    $('#put_url').focus();
  });
  $(document).keypress(function (e) {
    if ($("#put_modal").hasClass('in') && (e.keycode == 13 || e.which == 13)) {
      $('#put_modal .put_button').click();
    }
  });
}

function delete_modal_code() {
  $("#delete_file").val('');
  $('#delete_modal').modal('show');
  $('#delete_modal .delete_button').click(function () {
    $('#delete_modal').modal('hide');
    var to_delete = $("#delete_file").val();
    var uuid = $("#agent_uuid").val();
    C2deleteFile(uuid, to_delete);
  });
  $('#delete_modal').on('shown.bs.modal', function () {
    $('#delete_file').focus();
  });
  $(document).keypress(function (e) {
    if ($("#delete_modal").hasClass('in') && (e.keycode == 13 || e.which == 13)) {
      $('#delete_modal .delete_button').click();
    }
  });
}

function cycles_modal_code() {
  $("#setting_value").val('');
  $('#cycles_modal').modal('show');
  $('#cycles_modal .change_button').click(function () {
    $('#cycles_modal').modal('hide');
    var to_change_name = $("#setting_name").val();
    var to_change_value = $("#setting_value").val();
    var uuid = $("#agent_uuid").val();
    C2setValue(uuid, to_change_name, to_change_value);
  });
  $('#cycles_modal').on('shown.bs.modal', function () {
    $('#setting_name').focus();
  });
  $(document).keypress(function (e) {
    if ($("#cycles_modal").hasClass('in') && (e.keycode == 13 || e.which == 13)) {
      $('#cycles_modal .change_button').click();
    }
  });
}

function sleep_modal_code() {
  $("#sleep_seconds").val('');
  $('#sleep_modal').modal('show');
  $('#sleep_modal .sleep_button').click(function () {
    $('#sleep_modal').modal('hide');
    var to_sleep = $("#sleep_seconds").val();
    var uuid = $("#agent_uuid").val();
    C2sleepAgent(uuid, to_sleep);
  });
  $('#sleep_modal').on('shown.bs.modal', function () {
    $('#sleep_seconds').focus();
  });
}

function exit_modal_code() {
  $('#exit_modal').modal('show');
  $('#exit_modal .exit_button').click(function () {
    $('#exit_modal').modal('hide');
    var uuid = $("#agent_uuid").val();
    C2killAgent(uuid);
  });
}

function lock_modal_code() {
  $('#lock_modal').modal('show');
  $('#lock_modal .lock_button').click(function () {
    $('#lock_modal').modal('hide');
    var uuid = $("#agent_uuid").val();
    C2lockMachine(uuid);
  });
}

function reg_modal_code() {
  var uuid = $("#agent_uuid").val();
  C2lockMachine(uuid);
}

function confirm_modal_code() {
  $('#confirm_modal').modal('show');
  $('#confirm_modal .confirm_button').click(function () {
    $('#confirm_modal').modal('hide');
    var uuid = $("#agent_uuid").val();
    AgentDeleteAction([uuid]);
  });
}
