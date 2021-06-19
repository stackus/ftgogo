@command @delivery
Feature: Scheduling Deliveries

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

  Scenario: Deliveries can be scheduled with a courier
    When I schedule a delivery with:
    """
    {
      "OrderID": "a123",
      "ReadyBy": "2006-01-02T15:04:05Z"
    }
    """
    Then I expect the command to succeed

  Scenario: Scheduling a delivery sets the status to "SCHEDULED"
    When I schedule a delivery with:
    """
    {
      "OrderID": "a123",
      "ReadyBy": "2006-01-02T15:04:05Z"
    }
    """
    And I get the delivery with:
    """
    {
      "OrderID": "a123"
    }
    """
    Then I expect the request to succeed
    And the returned delivery status is:
    """
    SCHEDULED
    """

  Scenario: Couriers are given a two step delivery plan when scheduled with a delivery
    When I schedule a delivery with:
    """
    {
      "OrderID": "a123",
      "ReadyBy": "2006-01-02T15:04:05Z"
    }
    """
    And I get the delivery with:
    """
    {
      "OrderID": "a123"
    }
    """
    Then I get the courier with:
    """
    {
      "CourierID": "<AssignedCourierID>"
    }
    """
    And I expect the request to succeed
    And the returned courier matches:
    """
    {
      "CourierID": "<AssignedCourierID>",
      "Plan": [
        {
          "DeliveryID": "a123",
          "ActionType": "PICKUP",
          "Address": {
            "Street1": "123 Address St.",
            "City": "HomeTown",
            "State": "HomeState",
            "Zip": "12345"
          },
          "When": "2006-01-02T15:04:05Z"
        }, {
          "DeliveryID": "a123",
          "ActionType": "DROPOFF",
          "Address": {
            "Street1": "123 Address St.",
            "City": "HomeTown",
            "State": "HomeState",
            "Zip": "12345"
          },
          "When": "2006-01-02T15:34:05Z"
        }
      ],
      "Available": true
    }
    """

  Scenario: Unavailable couriers are not assigned to deliveries
    Given I create a courier with:
    """
    {
      "CourierID": "a123",
      "Available": false
    }
    """
    When I schedule a delivery with:
    """
    {
      "OrderID": "a123",
      "ReadyBy": "2006-01-02T15:04:05Z"
    }
    """
    And I get the delivery with:
    """
    {
      "OrderID": "a123"
    }
    """
    Then the returned delivery is not assigned to:
    """
    a123
    """
