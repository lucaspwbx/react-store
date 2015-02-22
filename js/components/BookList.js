var React = require('react');
var $ = require('jquery');
var Book = require('./Book');
var AddReviewForm = require('./AddReviewForm');

var BookList = React.createClass({
  handleNewReview: function(review) {
    console.log(review);
  },
  getInitialState: function() {
    return {
      books: []
    };
  },
  getBooks: function() {
    return $.ajax({
      url: 'http://localhost:8080/books',
      crossDomain: true,
      dataType: 'json',
    });
  },
  componentDidMount: function() {
    var _this = this;
    this.getBooks().then(function(result) {
      _this.setState({
        books: result
      });
    });
  },
  render: function() {
    var books = this.state.books.map(function(book, index) {
      var reviews = book.reviews ? book.reviews : [];

      return (
        <div>
          <Book key={book.id} title={book.title} language={book.language} pages={book.pages} reviews={reviews}/>
          <AddReviewForm bookId={book.id} newReview={this.handleNewReview} />
        </div>
      );
    }.bind(this));

    return (
      <div>
        <h1>Lista de livros</h1>
        {books}
      </div>
    );
  }
});

module.exports = BookList;
