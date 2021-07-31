@command @courier
Feature: Setting Courier Availability

  Scenario: Couriers can be created
    When I set the courier "Quick Courier" to be available
    Then I expect the command to succeed
