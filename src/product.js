import React from 'react';

class Product extends React.Component {
  handleClick() {
    let product = this.props.product;
    this.props.onAddProductToCart(product);
  }
  render() {
    return (
      <li>
        {this.props.product.name} -> R$ {this.props.product.price}
        <button onClick={this.handleClick.bind(this)}>
          Add to Cart
        </button>
      </li>
    )
  }
}

export default Product;
