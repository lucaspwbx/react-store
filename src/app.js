import React from 'react';
import Catalog from './catalog';
import NewProductForm from './productform';
import Cart from './cart';
import _ from 'underscore';

var products = [
  {id: 1, name: 'Shirt', price: 2},
  {id: 2, name: 'Jeans', price: 3}
];

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      products: products,
      cart: [],
    };
  }
  handleAddToCart(product) {
    var id = this.state.cart.length > 0 ? ++this.state.cart[this.state.cart.length-1].id : 1;
    product.id = id;
    var index = _.findIndex(this.state.cart, (item) => {
      return item.name == product.name;
    });
    if (index >= 0) {
      var newCart = this.state.cart;
      newCart[index].quantity++;
    } else {
      product.quantity = 1;
      var newCart = this.state.cart.concat(product);
    }
    this.setState({cart: newCart});
  }
  handleAddNewProduct(product) {
    product.id = ++this.state.products[this.state.products.length-1].id;
    let newProducts = this.state.products.concat(product);
    this.setState({products: newProducts});
  }
  render() {
    return (
      <div>
        <h1>App</h1>
        <Catalog onAddProductToCart={this.handleAddToCart.bind(this)} products={this.state.products}/>
        <NewProductForm onNewProduct={this.handleAddNewProduct.bind(this)}/>
        <Cart currentCart={this.state.cart}/>
      </div>
    )
  }
}

export default App;

React.render(<App/>, document.getElementById('app'));
