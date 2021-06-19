@command @delivery
Feature: Create Deliveries

  Background: Setup a restaurant
    Given I setup a restaurant with:
    """
    {
      "RestaurantID": "a123",
      "Name": "TestRestaurant",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """

  Scenario: Can create deliveries
    When I create a delivery with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "a123",
      "DeliveryAddress": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    Then I expect the command to succeed

  Scenario: Deliveries are created with a "PENDING" status
    When I create a delivery with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "a123",
      "DeliveryAddress": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    And I get the delivery with:
    """
    {
      "OrderID": "a123"
    }
    """
    Then I expect the command to succeed
    And the returned delivery status is:
    """
    PENDING
    """


  Scenario: Creating deliveries for restaurants that do not exist returns an error
    When I create a delivery with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "b456",
      "DeliveryAddress": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    restaurant not found
    """

  Scenario: Creating duplicate deliveries for an order returns an error
    Given I create a delivery with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "a123",
      "DeliveryAddress": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    When I create another delivery with:
    """
    {
      "OrderID": "a123",
      "RestaurantID": "a123",
      "DeliveryAddress": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    delivery already exists
    """
