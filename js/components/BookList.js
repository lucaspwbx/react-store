var React = require('react');
var $ = require('jquery');
var Book = require('./Book');

var BookList = React.createClass({
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
        <Book key={index} title={book.title} language={book.language} pages={book.pages} reviews={reviews}/>
      );
    });

    return (
      <div>
      Lista de livros
      {books}
      </div>
    );
  }
});

module.exports = BookList;
