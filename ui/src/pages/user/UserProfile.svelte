<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import UpdatePasswordForm from '../../components/user/UpdatePasswordForm.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import { warrior } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import ProfileForm from '../../components/user/ProfileForm.svelte';
  import CreateApiKey from '../../components/user/CreateApiKey.svelte';
  import DeleteConfirmation from '../../components/DeleteConfirmation.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;

  let warriorProfile = {};
  let apiKeys = [];
  let showApiKeyCreate = false;
  let showAccountDeletion = false;

  let updatePassword = false;

  const { ExternalAPIEnabled, LdapEnabled, HeaderAuthEnabled } = AppConfig;

  function toggleUpdatePassword() {
    updatePassword = !updatePassword;
    eventTag(
      'update_password_toggle',
      'engagement',
      `update: ${updatePassword}`,
    );
  }

  function getProfile() {
    xfetch(`/api/users/${$warrior.id}`)
      .then(res => res.json())
      .then(function (result) {
        warriorProfile = result.data;
      })
      .catch(function () {
        notifications.danger($LL.profileErrorRetrieving());
        eventTag('fetch_profile', 'engagement', 'failure');
      });
  }

  function updateWarriorProfile(p) {
    const body = {
      ...p,
    };

    xfetch(`/api/users/${$warrior.id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function () {
        warrior.update({
          id: warriorProfile.id,
          name: p.name,
          email: warriorProfile.email,
          rank: warriorProfile.rank,
          avatar: p.avatar,
          notificationsEnabled: p.notificationsEnabled,
          locale: p.locale,
        });

        notifications.success($LL.profileUpdateSuccess());
        eventTag('update_profile', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger($LL.profileErrorUpdating());
        eventTag('update_profile', 'engagement', 'failure');
      });
  }

  function updateWarriorPassword(password1, password2) {
    const body = {
      password1,
      password2,
    };

    xfetch('/api/auth/update-password', { body, method: 'PATCH' })
      .then(function () {
        notifications.success($LL.passwordUpdated(), 1500);
        toggleUpdatePassword();
        eventTag('update_password', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger($LL.passwordUpdateError());
        eventTag('update_password', 'engagement', 'failure');
      });
  }

  function getApiKeys() {
    xfetch(`/api/users/${$warrior.id}/apikeys`)
      .then(res => res.json())
      .then(function (result) {
        apiKeys = result.data;
      })
      .catch(function () {
        notifications.danger($LL.apiKeysErrorRetrieving());
        eventTag('fetch_profile_apikeys', 'engagement', 'failure');
      });
  }

  function deleteApiKey(apk) {
    return function () {
      xfetch(`/api/users/${$warrior.id}/apikeys/${apk}`, {
        method: 'DELETE',
      })
        .then(res => res.json())
        .then(function (result) {
          notifications.success($LL.apiKeyDeleteSuccess());
          apiKeys = result.data;
        })
        .catch(function () {
          notifications.danger($LL.apiKeyDeleteFailed());
        });
    };
  }

  function toggleApiKeyActiveStatus(apk, active) {
    return function () {
      const body = {
        active: !active,
      };

      xfetch(`/api/users/${$warrior.id}/apikeys/${apk}`, {
        body,
        method: 'PUT',
      })
        .then(res => res.json())
        .then(function (result) {
          notifications.success($LL.apiKeyUpdateSuccess());
          apiKeys = result.data;
        })
        .catch(function () {
          notifications.danger($LL.apiKeyUpdateFailed());
        });
    };
  }

  function toggleCreateApiKey() {
    showApiKeyCreate = !showApiKeyCreate;
  }

  function handleDeleteAccount() {
    xfetch(`/api/users/${$warrior.id}`, { method: 'DELETE' })
      .then(function () {
        warrior.delete();

        eventTag('delete_warrior', 'engagement', 'success');

        router.route(appRoutes.landing);
      })
      .catch(function () {
        notifications.danger($LL.profileDeleteError());
        eventTag('delete_warrior', 'engagement', 'failure');
      });
  }

  function toggleDeleteAccount() {
    showAccountDeletion = !showAccountDeletion;
  }

  onMount(() => {
    if (!$warrior.id) {
      router.route(appRoutes.login);
      return;
    }

    getProfile();
    if (ExternalAPIEnabled) {
      getApiKeys();
    }
  });
</script>

<svelte:head>
  <title>{$LL.profileTitle()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight dark:text-white"
  >
    {$LL.profileTitle()}
  </h1>

  <div class="flex justify-center flex-wrap">
    <div class="w-full md:w-1/2 lg:w-1/3">
      {#if !updatePassword}
        <div
          class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 md:p-6 mb-4"
        >
          <ProfileForm
            profile="{warriorProfile}"
            handleUpdate="{updateWarriorProfile}"
            toggleUpdatePassword="{toggleUpdatePassword}"
            xfetch="{xfetch}"
            notifications="{notifications}"
            eventTag="{eventTag}"
            ldapEnabled="{LdapEnabled}"
            headerAuthEnabled="{HeaderAuthEnabled}"
          />
        </div>
      {/if}

      {#if updatePassword}
        <div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4">
          <div
            class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
        md:leading-tight text-center dark:text-white"
          >
            {$LL.updatePassword()}
          </div>

          <UpdatePasswordForm
            handleUpdate="{updateWarriorPassword}"
            toggleForm="{toggleUpdatePassword}"
            notifications="{notifications}"
          />
        </div>
      {/if}
    </div>

    <div class="w-full md:w-1/2 lg:w-2/3">
      {#if ExternalAPIEnabled}
        <div class="ms-8">
          <div class="flex w-full">
            <div class="flex-1">
              <h2
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
              >
                {$LL.apiKeys()}
              </h2>
            </div>
            <div class="flex-1">
              <div class="text-right">
                <HollowButton
                  href="/swagger/index.html"
                  options="{{ target: '_blank' }}"
                  color="blue"
                >
                  {$LL.apiDocumentation()}
                </HollowButton>
                <HollowButton
                  onClick="{toggleCreateApiKey}"
                  testid="apikey-create"
                >
                  {$LL.apiKeyCreateButton()}
                </HollowButton>
              </div>
            </div>
          </div>

          <div class="flex flex-col">
            <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
              <div
                class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
              >
                <div
                  class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
                >
                  <table
                    class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                  >
                    <thead class="bg-gray-50 dark:bg-gray-800">
                      <tr>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.name()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.apiKeyPrefix()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.active()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.lastUpdated()}
                        </th>
                        <th scope="col" class="relative px-6 py-3">
                          <span class="sr-only">{$LL.actions()}</span>
                        </th>
                      </tr>
                    </thead>
                    <tbody
                      class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                    >
                      {#each apiKeys as apk, i}
                        <tr
                          class:bg-slate-100="{i % 2 !== 0}"
                          class:dark:bg-gray-800="{i % 2 !== 0}"
                          data-testid="apikey"
                          data-apikeyid="{apk.id}"
                        >
                          <td
                            class="px-6 py-4 whitespace-nowrap"
                            data-testid="apikey-name">{apk.name}</td
                          >
                          <td
                            class="px-6 py-4 whitespace-nowrap"
                            data-testid="apikey-prefix"
                          >
                            {apk.prefix}
                          </td>
                          <td
                            class="px-6 py-4 whitespace-nowrap"
                            data-testid="apikey-active"
                            data-active="{apk.active}"
                          >
                            {#if apk.active}
                              <span class="text-green-600"><CheckIcon /></span>
                            {/if}
                          </td>
                          <td class="px-6 py-4 whitespace-nowrap">
                            {new Date(apk.updatedDate).toLocaleString()}
                          </td>
                          <td
                            class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                          >
                            <HollowButton
                              onClick="{toggleApiKeyActiveStatus(
                                apk.id,
                                apk.active,
                              )}"
                              testid="apikey-activetoggle"
                            >
                              {#if !apk.active}
                                {$LL.activate()}
                              {:else}
                                {$LL.deactivate()}
                              {/if}
                            </HollowButton>
                            <HollowButton
                              color="red"
                              onClick="{deleteApiKey(apk.id)}"
                              testid="apikey-delete"
                            >
                              {$LL.delete()}
                            </HollowButton>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    </div>

    {#if !LdapEnabled && !HeaderAuthEnabled}
      <div class="w-full text-center mt-8">
        <HollowButton onClick="{toggleDeleteAccount}" color="red">
          {$LL.deleteAccount()}
        </HollowButton>
      </div>
    {/if}
  </div>
  {#if showApiKeyCreate}
    <CreateApiKey
      toggleCreateApiKey="{toggleCreateApiKey}"
      handleApiKeyCreate="{getApiKeys}"
      notifications="{notifications}"
      xfetch="{xfetch}"
      eventTag="{eventTag}"
    />
  {/if}

  {#if showAccountDeletion}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteAccount}"
      handleDelete="{handleDeleteAccount}"
      confirmText="{$LL.deleteAccountWarningStatement()}"
      confirmBtnText="{$LL.deleteConfirmButton()}"
    />
  {/if}
</PageLayout>
