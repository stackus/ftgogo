@command @order @create
Feature: Order Creation

  Background: Initialize resources
    Given I have initialized the restaurant "Best Foods"

  Scenario: Can create new orders
    When I submit an order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed

  Scenario: Cannot create new orders at non-existent restaurants
    When I submit an order to "Other Foods" from "Able Anders"
    Then I expect the command to fail
