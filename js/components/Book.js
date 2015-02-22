var React = require('react');
var Review = require('./Review');

var Book = React.createClass({
  render: function() {
    var reviews = '';
    if (this.props.reviews.length > 0) {
      reviews = this.props.reviews.map(function(review, index) {
        return (
          <Review key={index} description={review.description} name={review.name} />
        );
      });
    } else {
      reviews = 'No reviews';
    }
    return (
      <div>
        <p>Title: {this.props.title}</p>
        <p>Language: {this.props.language}</p>
        <p>Pages: {this.props.pages}</p>
        <ul>
          {reviews}
        </ul>
        <hr/>
      </div>
    );
  }
});

module.exports = Book;
