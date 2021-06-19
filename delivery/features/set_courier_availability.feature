@command @courier
Feature: Setting Courier Availability

  Scenario: Couriers can be created
    When I set a couriers availability with:
    """
    {
      "CourierID": "a123",
      "Available": true
    }
    """
    Then I expect the command to succeed
