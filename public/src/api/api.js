var auth = require("../app/auth.js");
var navigate = require('react-mini-router').navigate;

module.exports = function(url, callback, errorCallback) {
    $.ajax({
        url: url,
        dataType: 'json',
        beforeSend: function (request) {
            request.setRequestHeader("X-Auth-Token", auth.getUsername() + ":" + auth.getToken());
        },
        cache: false,
        success: function(data) {
            callback(data);
        },
        error: function(xhr, status, err) {
            if (status == 401) {
                navigate("/login");
                return;
            }
            errorCallback(xhr, status, err);
        }
    });
};
