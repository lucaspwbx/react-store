import React from 'react';
import Catalog from './catalog';
import CartProduct from './cartproduct';
import _ from 'underscore';

class Cart extends React.Component {
  render() {
    var bla = _.map(this.props.currentCart, (item) => {
      return item.price * item.quantity;
    })
    var total = _.reduce(bla, (memo, num) => {
      return memo + num;
    }, 0);
    var cart = this.props.currentCart.map((item, i) => {
      return <CartProduct product={item} key={i}/>;
    });
    return (
      <div>
        {cart}
        <p>Total carrinho: R$ {total}</p>
      </div>
    )
  }
}

export default Cart;
