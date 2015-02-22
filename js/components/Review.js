var React = require('react');

var Review = React.createClass({
  render: function() {
    return (
      <div>
        <p>{this.props.description}</p>
        <p>User: {this.props.name}</p>
      </div>
    );
  }
});

module.exports = Review;
