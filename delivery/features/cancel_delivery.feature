@command @delivery
Feature: Cancel Deliveries

  Background: Setup resources
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
    And I setup a delivery with:
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

  Scenario: Cancel existing deliveries
    When I cancel the delivery with:
    """
    {
      "OrderID": "a123"
    }
    """
    Then I expect the command to succeed

  Scenario: Canceling deliveries that do not exist returns an error
    When I cancel a delivery with:
    """
    {
      "OrderID": "b456"
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    delivery not found
    """
