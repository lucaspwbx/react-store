var React = require('react');

var Book = React.createClass({
  render: function() {
    var reviews = '';
    if (this.props.reviews.length > 0) {
      reviews = 'Reviews: ' + this.props.reviews;
    } else {
      reviews = 'No reviews';
    }
    return (
      <div>
      Title: {this.props.title}
      Language: {this.props.language}
      Pages: {this.props.pages}
      {reviews}
      </div>
    );
  }
});

module.exports = Book;
