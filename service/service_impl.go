package service

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gomarket-loyalty/constants"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository"
	"strings"
	"time"
)

func NewProductService(Repository *repository.Repository) Service {
	return &serviceImpl{
		repository: *Repository,
		client:     ClientJSON{},
	}
}

type serviceImpl struct {
	repository repository.Repository
	client     Client
}

func (service *serviceImpl) Base() string {
	return "Hello World"
}

func (service *serviceImpl) Create(ctx context.Context, request model.RegisterRequest) error {
	err := service.ValidateDataRegister(request)
	if err != nil {
		return err
	}

	user := model.User{
		Login: request.Login,
	}

	err = service.repository.SetUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *serviceImpl) AddMechanic(ctx context.Context, bonus model.Mechanic) error {
	if bonus.RewardType != constants.Points && bonus.RewardType != constants.Percentage {
		return exception.ErrEnabledData
	}
	if bonus.Reward <= 0 {
		return exception.ErrEnabledData

	}
	err := service.repository.AddMechanic(ctx, bonus)
	if err != nil {
		return err
	}
	return nil

}

// ValidateDataRegister validates the user data for registration.
// It checks if the login is between 3 and 15 characters long,
// and if the password only contains alphanumeric characters.
// If any validation fails, it returns the ErrEnabledData error.
func (service *serviceImpl) ValidateDataRegister(user model.RegisterRequest) error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Login, validation.Required, validation.Length(3, 15)),
	)

	if err != nil {
		return exception.ErrEnabledData
	}
	return nil
}

// CreateOrder creates an order for a client with the specified clientID, orderID, and items.
func (service *serviceImpl) CreateOrder(ctx context.Context, clientID string, orderID string, order model.Items) error {
	// Initialize variables
	var bonus int
	var goodID []int

	// Create a map to store the count of each item ID
	idCount := make(map[int]int)

	// Iterate through each item in the order
	for _, item := range order.Items {
		// Skip items with invalid count or price
		if item.Count <= 0 || item.Price <= 0 {
			continue
		}

		// Add the item ID and count to the map
		idCount[item.Id] = item.Count

		// Append the item ID to the goodID slice
		goodID = append(goodID, item.Id)
	}

	// Create a request to get the name of the items with the specified IDs
	reqItems := model.RequestNameItems{
		Offset: 0,
		Limit:  1000,
		Query:  model.Query{Ids: goodID},
	}

	// Send the request and get the response
	var nameItems model.ResponseNameItems
	resItems, err := service.client.JSONRequest(reqItems, &nameItems, constants.URLGETNameItems)
	if err != nil {
		return fmt.Errorf("error getting name items: %w", err)
	}

	// Get all mechanics from the repository
	mechanics, err := service.repository.GetAllMechanics(ctx)
	if err != nil {
		return fmt.Errorf("error getting mechanics: %w", err)
	}

	switch v := resItems.(type) {
	case model.ResponseNameItems:
		// Iterate through each item in the response
		for _, item := range v.Items {
			// Iterate through each mechanic
			for _, mechanic := range mechanics {
				// Check if the item name contains the mechanic's match string
				if strings.Contains(item.Name, mechanic.Match) {
					// Create an order item with the item ID, price, and count
					order := model.Item{
						Id:    item.ID,
						Price: item.Price,
						Count: idCount[item.ID],
					}

					// Add the bonus for the mechanic to the total bonus
					bonus += service.AddBonus(mechanic, order)
				}
			}
		}
	default:
		return fmt.Errorf("invalid response type: %T", v)
	}

	// Create an order transaction with the order ID and bonus
	transaction := model.Order{
		User:  clientID,
		Order: orderID,
		Bonus: bonus,
		Time:  time.Now().Format(constants.FormatTime),
	}

	// Create the order transaction in the repository
	err = service.repository.CreateOrder(ctx, transaction)
	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	// Update the bonus for the client in the repository
	err = service.repository.UpdateBonusUser(ctx, clientID, bonus)
	if err != nil {
		return fmt.Errorf("error updating bonus: %w", err)
	}

	return nil
}

func (service *serviceImpl) AddBonus(mechanic model.Mechanic, item model.Item) int {
	var oneBonus int
	switch mechanic.RewardType {
	case constants.Points:
		oneBonus += mechanic.Reward
	case constants.Percentage:
		oneBonus = (item.Price / 100) * mechanic.Reward
	}
	return oneBonus * item.Count
}

func (service *serviceImpl) GetInfoOrders(ctx context.Context, id string) ([]model.Order, error) {

	orders, err := service.repository.GetInfoOrders(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, exception.ErrNotFound
	}
	return orders, nil

}
