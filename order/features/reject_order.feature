@command @order @reject
Feature: Order Rejection

  Background: Setup Resources
    Given I have initialized the restaurant "Best Foods"
    And I have submitted an order to "Best Foods" from "Able Anders"

  Scenario: Pending orders may be rejected
    Given the order to "Best Foods" from "Able Anders" is "ApprovalPending"
    When I reject the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And the order to "Best Foods" from "Able Anders" is "Rejected"

  Scenario: Approved orders cannot be rejected
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    When I reject the order to "Best Foods" from "Able Anders"
    Then I expect the command to fail
    And the returned error message is "order state does not allow action"
