var Menu = require('../app/menu.js');

module.exports = React.createClass({
    render: function () {
        return (
            <div className="pusher">
                <Menu />
                <div className="app-content">
                </div>
            </div>
        )
    }
});
