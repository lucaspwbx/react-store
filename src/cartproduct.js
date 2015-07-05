import React from 'react';

class CartProduct extends React.Component {
  render() {
    return (
      <p>{this.props.product.name} -> R$ {this.props.product.price}, Quantity: {this.props.product.quantity} = R$ {this.props.product.quantity * this.props.product.price}</p>
    )
  }
}

export default CartProduct;
