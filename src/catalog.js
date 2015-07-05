import React from 'react';
import Product from './product';

class Catalog extends React.Component {
  handleAddToCart(product) {
    this.props.onAddProductToCart(product);
  }
  render() {
    var products = this.props.products.map((product,i) => {
      return <Product onAddProductToCart={this.handleAddToCart.bind(this)} product={product} key={i}/>;
    });
    return (
      <div>
        Catalogo
        <ul>
        {products}
        </ul>
      </div>
    )
  }
}

export default Catalog;
