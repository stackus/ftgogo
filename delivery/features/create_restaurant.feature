@command @restaurant
Feature: Create Restaurants

  Scenario: Can create new restaurants
    When I create a restaurant with:
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
    Then I expect the command to succeed

  Scenario: Creating duplicate restaurants returns an error
    Given I create a restaurant with:
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
    When I create another restaurant with:
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
    Then I expect the command to fail
    And the returned error message is:
    """
    restaurant already exists
    """
