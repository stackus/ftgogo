@command @order @cancel
Feature: Order Cancellation

  Background: Setup Resources
    Given I have initialized the restaurant "Best Foods"
    And I have submitted an order to "Best Foods" from "Able Anders"

  Scenario: Approved orders may be cancelled
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    When I begin to cancel the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "CancelPending"

  Scenario: Rejected orders cannot be cancelled
    Given I have rejected the order to "Best Foods" from "Able Anders"
    When I begin to cancel the order to "Best Foods" from "Able Anders"
    Then I expect the command to fail
    And the returned error message is "order state does not allow action"

  Scenario: Can confirm cancellation for orders pending cancellation
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    And I have begun to cancel the order to "Best Foods" from "Able Anders"
    When I confirm cancelling the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "Cancelled"

  Scenario: Can undo cancellation for orders pending cancellation
    Given I have approved the order to "Best Foods" from "Able Anders" with ticket "T123"
    And I have begun to cancel the order to "Best Foods" from "Able Anders"
    When I undo cancelling the order to "Best Foods" from "Able Anders"
    Then I expect the command to succeed
    And expect the order to "Best Foods" from "Able Anders" is "Approved"
