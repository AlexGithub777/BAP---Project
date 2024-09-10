# BAP Testing
## Test 1: Login

**Component:** Authentication

**Prerequisites**: Standard login

**Steps**
* User opens app
* User inputs required information in fields (ie Email in email field, password in password field)
* User clicks login

**Test Results**
<table>
    <tr>
        <th>Input</th>
        <th>Expected Result</th>
        <th>Actual Result</th>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: "Password1!"</td>
        <td>User successfully logs in and can see a map of EIT's Campus</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: <BLANK><br>Password: "Password1!"</td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: <BLANK></td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: <BLANK><br>Password: <BLANK></td>
        <td>Login Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
</table>

___

## Test 2: Add Emergency Device

**Component:** Dashboard 

**Prerequisites**: Admin account

**Steps**
* User clicks on add device
* User inputs relevant information pertaining to the fire extinguisher
* User clicks add

**Test Results**
<table>
    <tr>
        <th>Input</th>
        <th>Expected Result</th>
        <th>Actual Result</th>
    </tr>
    <tr>
        <td>Type: Fire Extinguisher<br>Extinguisher Type: CO2<br>Building: A<br>Room: A1<br>Serial Number: SN00001<br>Description: By charging station<br>Manufactured: 6/7/24<br>Last inspected: 8/9/24<br>Status: Active</td>
        <td>User successfully adds a fire extinguisher type device into the system</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Type: Fire Blanket<br>Building: B<br>Room: B1<br>Description: On wall outside classroom<br>Manufactured: 6/8/20</td>
        <td>User successfully adds a fire blanket type device into the system</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Type: Med-Kit<br>Building: C<br>Room: C1<br>Serial Number: MD00001<br>Description: In a cupboard</td>
        <td>User successfully adds a med kit into the system</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Any blank information</td>
        <td>Adding device fails due to missing information</td>
        <td>UNTESTED</td>
    </tr>
</table>

___

## Test 3: User account to Admin account

**Component**: Administrator

**Prerequisites**: Admin account (isAdmin=True)

**Steps**
* User opens user list
* User finds other user to change into admin (or searches for)
* User clicks on them, then on a drop down field saying “Set Admin”
* User can repeat steps as many times as required

**Test data**: None

**Expected Results**: User can update a user into an Admin

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

___ 

## Test 4: View Emergency Device

**Component**: Dashboard

**Prerequisites**: User account

**Steps**
* User finds an extinguisher to open (map)
* User clicks on fire extinguisher

**Test data**: None

**Expected Results**: User can view a fire extinguishers info

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

___ 

## Test 5: Sign up

**Component**: Authentication

**Prerequisites**: None

**Steps**
* User opens app and clicks on "sign up"
* User fills in fields with relevant information
* User clicks sign up

**Test Results**
<table>
    <tr>
        <th>Input</th>
        <th>Expected Result</th>
        <th>Actual Result</th>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: "Password1!"<br>Email:"user1email@gmail.com"</td>
        <td>User successfully signs up and gets redirected to log in page</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: "Password1!"<br>Email:<BLANK></td>
        <td>Sign up Fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: "user1"<br>Password: <BLANK><br>Email:"user1email@gmail.com"</td>
        <td>Sign up fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Username: <BLANK><br>Password: "Password1!"<br>Email:"user1email@gmail.com"</td>
        <td>Sign up fails due to missing input/s</td>
        <td>UNTESTED</td>
    </tr>
</table>

___ 

## Test 6: Deleting account

**Component**: Dashboard

**Prerequisites**: User account

**Steps**
* User logs into account
* User clicks on account settings in app (From home: Settings -> Account Settings)
* User clicks delete account

**Test data**: None

**Expected Results**: User account is successfully deleted and logs user out

**Actual Results**: UNTESTED

**Test Status**: UNTESTED

## Test 7: Device Inspection

**Component**: Dashboard

**Prerequisites**: User account

**Steps**
* User logs into account
* User clicks on account settings in app (From home: Settings -> Account Settings)
* User clicks delete account

**Test data**: None
<table>
    <tr>
        <th>Input</th>
        <th>Expected Result</th>
        <th>Actual Result</th>
    </tr>
    <tr>
        <td>Inspection Date: 8/9/24<br>Notes: Small dent<br>Is conspicuous: true<br>Operation instructions clear: true<br>Is accessible: true<br>Is in correct location: true<br>Has maintenance tag: true<br>Is sign visible: true<br>Is external damage present: true<br>Is anti-tamper device intact: true<br>is charge gauge normal: true<br>is support bracket secure: true<br>Is replaced: false<br>Is Work order required: false<br>Maintenance records complete: true</td>
        <td>Device passes inspection</td>
        <td>UNTESTED</td>
    </tr>
    <tr>
        <td>Inspection Date: 8/9/24<br>Notes: Gone<br>Is conspicuous: false<br>Operation instructions clear: true<br>Is accessible: false<br>Is in correct location: false<br>Has maintenance tag: false<br>Is sign visible: true<br>Is external damage present: true<br>Is anti-tamper device intact: false<br>is charge gauge normal: false<br>is support bracket secure: true<br>Is replaced: false<br>Is Work order required: true<br>Maintenance records complete: false</td>
        <td>Device fails inspection</td>
        <td>UNTESTED</td>
    </tr>
</table>

___

## Test 8: View Notifications

**Component**: Notification

**Prerequisites**: User account

**Steps**
* User logs into account
* User clicks on notifications to check if they have anything that needs to be updated

**Test data**: None

**Expected Results**: User can see notifications if anything new has occurred

**Actual Results**: UNTESTED

**Test Status**: UNTESTED
