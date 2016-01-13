var navigate = require('react-mini-router').navigate;
var auth = require('./auth.js');
var UserMenu = require('./usermenu.js');

module.exports = React.createClass({
    getInitialState: function() {
        return {
            username: auth.getUsername(),
            isLoggedIn: auth.isLoggedIn()
        }
    },
    componentDidMount() {
        this.setState({
            username: auth.getUsername(),
            isLoggedIn: auth.isLoggedIn()
        });
    },
    click: function(path) {
        navigate(path);
    },
    render: function () {
        return (
            <div className="ui left fixed vertical inverted menu">
                <div className="item">
                    <a onClick={this.click.bind(this, "/")}>
                    <img className="ui avatar tiny image" src="./assets/images/quorra_logo.png"/>
                    <span className="logo-text">Quorra</span>
                    </a>
                </div>
                <a onClick={this.click.bind(this, "/")}className="item">Dashboard</a>
                <div className="item">
                    <div className="header">Zones</div>
                    <div className="menu">
                        <a className="item">First</a>
                        <a className="item">Second</a>
                        <a className="item">Basement</a>
                        <a className="item">Office</a>
                        <a className="item">Exterior</a>
                    </div>
                </div>

                {
                    this.state.isLoggedIn ? (
                        <UserMenu />
                    ) : (
                        <div></div>
                    )
                }
                <a onClick={this.click.bind(this, "/login")}className="item">Login</a>
                <a onClick={this.click.bind(this, "/help")}className="item">Help</a>
            </div>
        )
    }
});
