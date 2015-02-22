var React = require('react');
var BookList = require('./BookList');

var BookApp = React.createClass({
  render: function() {
    return (
      <div>
      <BookList/>
      </div>
    );
  }
});

module.exports = BookApp;
