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
        <p>Title: {this.props.title}</p>
        <p>Language: {this.props.language}</p>
        <p>Pages: {this.props.pages}</p>
        <p>{reviews}</p>
        <hr/>
      </div>
    );
  }
});

module.exports = Book;
