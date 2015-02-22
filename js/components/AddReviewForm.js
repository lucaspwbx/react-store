var React = require('react');

var AddReviewForm = React.createClass({
  handleSubmit: function(e) {
    e.preventDefault();
    var description = this.refs.description.getDOMNode().value;
    var name = this.refs.name.getDOMNode().value;
    this.props.newReview({id: this.props.bookId, description: description, name: name});
  },
  render: function() {
    return (
      <form onSubmit={this.handleSubmit}>
        <input type='text' ref='description'/>
        <input type='text' ref='name'/>
        <input type='submit' value='Save Review'/>
      </form>
    );
  }
});

module.exports = AddReviewForm;
