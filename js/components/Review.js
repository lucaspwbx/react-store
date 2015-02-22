var React = require('react');

var Review = React.createClass({
  render: function() {
    return (
      <li>
        <p>{this.props.description}</p>
        <p>User: {this.props.name}</p>
      </li>
    );
  }
});

module.exports = Review;
