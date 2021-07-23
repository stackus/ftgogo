@query @courier
Feature: Get Couriers

  Background: Setup a Courier
    Given a courier exists named "Quick Courier"

  Scenario: Can get couriers
    When I get the courier named "Quick Courier"
    Then I expect the request to succeed
    And the returned courier is available

  Scenario: Can get unavailable couriers
    Given I set the courier "Quick Courier" to be unavailable
    When I get the courier named "Quick Courier"
    Then I expect the request to succeed
    And the returned courier is not available
