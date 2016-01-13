var Login = require("./login.js");
var Dashboard = require("../dashboard/dashboard.js");
var NotFound = require("./notfound.js");

var RouterMixin = require('react-mini-router').RouterMixin;

var App = React.createClass({
    mixins: [RouterMixin],
    routes: {
        '/': 'home',
        '/login': 'login',
    },
    render: function() {
        return this.renderCurrentRoute();
    },
    home: function() {
        return (
            <Dashboard />
        )
    },
    login: function() {
        return (
            <Login />
        )
    },
    notFound: function(path) {
        this.state.path = path;
        return (
            <NotFound path={this.state.path} />
        )
    }
});

React.render(<App />, document.querySelector('body'));
