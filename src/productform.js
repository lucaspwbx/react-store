import React from 'react';

class NewProductForm extends React.Component {
  handleSubmit(e) {
    e.preventDefault();
    var name = this.refs.name.getDOMNode().value;
    var price = this.refs.price.getDOMNode().value;
    this.props.onNewProduct({name: name, price: price});
  }
  render() {
    return (
      <div>
        <form onSubmit={this.handleSubmit.bind(this)}>
          <label htmlFor="name">Nome: </label>
          <input type="text" ref="name" />
          <label htmlFor="price">Preco: </label>
          <input type="text" ref="price" />
          <button type="submit">Clique</button>
        </form>
      </div>
    )
  }
}

export default NewProductForm;
