package database

import "errors"

var {
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodePrducts = errors.New("can't find the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this product from the cart")
	ErrCantGetItem = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartIem = errors.New("cannot update the purchase")
}

func AddProductToCart(){

}

func RemoveCartItem(){

}

func BuyItemFromCart(){

}

func InstantBuyer(){

}