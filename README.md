# Github Dependabot Alert Report

## Quickstart

1. Install (one-time setup)
    1. Download or clone this repository
    1. `cd github-reports`
    1. `go build .`
1. [Create credentials](#create-credentials) (one-time setup)
1. Run! e.g.:

    ```sh
    ./github-dependabot-report
    --privateKeyFile=dependabot-report.2024-01-24.private-key.pem \
    --appID=807018 \
    --appInstallID=46514769  \
    --excludeNonProd \
    --excludeArchived \
    --outputFile=report.html
    ```

I will provide a compiled version when I figure out how to get around the MacOS
notarization issue.

## Create credentials

### Create a Github App

[Github recommends the use of Github Apps to authenticate with their REST APIs](https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api?apiVersion=2022-11-28#authenticating-with-a-token-generated-by-an-app),
and in fact, [the ones we are using here](https://docs.github.com/en/enterprise-cloud@latest/rest/dependabot/alerts?apiVersion=2022-11-28#list-dependabot-alerts-for-a-repository) require it, so
the first thing you need to do is create a Github App.
Just follow [these instructions](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app), with the following additions:

* You are creating a _GitHub App owned by a personal account_ (step 2)
* For the _Homepage URL_ you can really use anything you want (e.g. `https://gihub.com`)
* Uncheck _Active_ under __Webhook__
* In permissions, set the Repository permissions for __Dependabot alerts__ to Read-only
* Set __Where can this GitHub App be installed?__ to __Any account__

> You do not need to provide any other information.

Now click on "Create Github App".

Next we need to obtain a Private Key from this Github App which we will use
later during authentication. To do this, click on Edit next to the Github App
you just created, __Private keys__ section and click on
_Generate a private key_. This will download the private key to your machine
(as a `.pem` file); remember where it is stored as you will need this later.

While you are on this page, also take note of your __App ID__ which is visible
at the top of this page, you will need this later as well.

## Install the Github App in the GLCP Organization

Follow [these instructions](https://docs.github.com/en/apps/using-github-apps/installing-your-own-github-app),
with the following additions:

* in step 7, make sure you select "glcp" to install the app in the GLCP
Organization.
* In step 8, select **All repositories**

Now click on install.

You will need a final bit of information from the installed GitHub App, namely
the Github App installation ID. The only way I have been able to retrieve this
is by going to the _Settings_ page of the installed Github App (click on the gear
icon next to "Installed"). If you now look at the URL in the browser, you
should see `https://github.com/apps/dependabot-report/installations/12345678`:
`12345678` is the Github App installation ID (obviously your number will
be different).
