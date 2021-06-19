@query @courier
Feature: Get Couriers

  Background: Setup a Courier
    Given I setup a courier with:
    """
    {
      "CourierID": "a123",
      "Available": true
    }
    """

  Scenario: Can get couriers
    When I get the courier with:
    """
    {
      "CourierID": "a123"
    }
    """
    Then I expect the request to succeed
    And the returned courier is available

  Scenario: Can get unavailable couriers
    Given I set the couriers availability with:
    """
    {
      "CourierID": "a123",
      "Available": false
    }
    """
    When I get the courier with:
    """
    {
      "CourierID": "a123"
    }
    """
    Then I expect the request to succeed
    And the returned courier is not available
