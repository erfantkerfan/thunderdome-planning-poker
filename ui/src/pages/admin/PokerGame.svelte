<script lang="ts">
  import { onMount } from 'svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import { warrior } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import DeleteConfirmation from '../../components/DeleteConfirmation.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  // export let eventTag
  export let battleId;

  let showDeleteBattle = false;

  let battle = {
    name: '',
    votingLocked: false,
    autoFinishVoting: false,
    activePlanId: '',
    pointValuesAllowed: [],
    pointAverageRounding: '',
    users: [],
    plans: [],
    createdDate: '',
    updatedDate: '',
  };

  function getBattle() {
    xfetch(`/api/battles/${battleId}`)
      .then(res => res.json())
      .then(function (result) {
        battle = result.data;
      })
      .catch(function () {
        notifications.danger(
          $LL.getBattleError({ friendly: AppConfig.FriendlyUIVerbs }),
        );
      });
  }

  function deleteBattle() {
    xfetch(`/api/battles/${battleId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        router.route(appRoutes.adminPokerGames);
      })
      .catch(function () {
        notifications.danger(
          $LL.deleteBattleError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
      });
  }

  const toggleDeleteBattle = () => {
    showDeleteBattle = !showDeleteBattle;
  };

  onMount(() => {
    if (!$warrior.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($warrior)) {
      router.route(appRoutes.landing);
      return;
    }

    getBattle();
  });
</script>

<svelte:head>
  <title
    >{$LL.battles({ friendly: AppConfig.FriendlyUIVerbs })}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="battles">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
    >
      {battle.name}
    </h1>
  </div>

  <div class="w-full">
    <div class="p-4 md:p-6">
      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.votingLocked()}
          </HeadCol>
          <HeadCol>
            {$LL.autoFinishVoting()}
          </HeadCol>
          <HeadCol>
            {$LL.pointValuesAllowed()}
          </HeadCol>
          <HeadCol>
            {$LL.pointAverageRounding()}
          </HeadCol>
          <HeadCol>
            {$LL.dateCreated()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          <TableRow itemIndex="{0}">
            <RowCol>
              {#if battle.votingLocked}
                <span class="text-green-600"><CheckIcon /></span>
              {/if}
            </RowCol>
            <RowCol>
              {#if battle.autoFinishVoting}
                <span class="text-green-600"><CheckIcon /></span>
              {/if}
            </RowCol>
            <RowCol>
              {battle.pointValuesAllowed.join(', ')}
            </RowCol>
            <RowCol>
              {battle.pointAverageRounding}
            </RowCol>
            <RowCol>
              {new Date(battle.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(battle.updatedDate).toLocaleString()}
            </RowCol>
          </TableRow>
        </tbody>
      </Table>
    </div>
    <div class="p-4 md:p-6">
      <h3
        class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
      >
        {$LL.users()}
      </h3>

      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.abandoned()}
          </HeadCol>
          <HeadCol>
            {$LL.spectator()}
          </HeadCol>
          <HeadCol>
            {$LL.leader()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each battle.users as user, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId="{user.id}"
                      avatar="{user.avatar}"
                      gravatarHash="{user.gravatarHash}"
                      width="48"
                      class="h-10 w-10 rounded-full"
                    />
                  </div>
                  <div class="ms-4">
                    <div class="text-sm font-medium text-gray-900">
                      <a
                        href="{appRoutes.adminUsers}/{user.id}"
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        >{user.name}</a
                      >
                      {#if user.country}
                        &nbsp;
                        <CountryFlag
                          country="{user.country}"
                          additionalClass="inline-block"
                          width="32"
                          height="24"
                        />
                      {/if}
                    </div>
                  </div>
                </div>
              </RowCol>
              <RowCol>
                {user.rank}
              </RowCol>
              <RowCol>
                {#if user.active}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                {#if user.abandoned}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                {#if user.spectator}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                {#if battle.leaders.includes(user.id)}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </div>

    <div class="p-4 md:p-6">
      <h3
        class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
      >
        {$LL.plans({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>

      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.planReferenceId()}
          </HeadCol>
          <HeadCol>
            {$LL.voteCount()}
          </HeadCol>
          <HeadCol>
            {$LL.points()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.skipped()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each battle.plans as plan, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                {plan.name}
              </RowCol>
              <RowCol>
                {plan.type}
              </RowCol>
              <RowCol>
                {plan.referenceId}
              </RowCol>
              <RowCol>
                {plan.votes.length}
              </RowCol>
              <RowCol>
                {plan.points}
              </RowCol>
              <RowCol>
                {#if plan.active}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                {#if plan.skipped}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>

      <div class="text-center mt-4">
        <HollowButton
          color="red"
          onClick="{toggleDeleteBattle}"
          testid="battle-delete"
        >
          {$LL.battleDelete({ friendly: AppConfig.FriendlyUIVerbs })}
        </HollowButton>
      </div>

      {#if showDeleteBattle}
        <DeleteConfirmation
          toggleDelete="{toggleDeleteBattle}"
          handleDelete="{deleteBattle}"
          confirmText="{$LL.deleteBattleConfirmText({
            friendly: AppConfig.FriendlyUIVerbs,
          })}"
          confirmBtnText="{$LL.deleteBattle({
            friendly: AppConfig.FriendlyUIVerbs,
          })}"
        />
      {/if}
    </div>
  </div>
</AdminPageLayout>
