<script lang="ts">
	import pb from '$lib/pocketbase';
	import { goto } from '$app/navigation';

	type ViewMode = 'cps' | 'characters';
	let mode = $state<ViewMode>('cps');
	let items = $state<any[]>([]);
	let isLoading = $state(true);
	let deleteModal = $state() as HTMLDialogElement;
	let itemToDelete = $state<null | any>(null);

	const handleDelete = (item: any) => {
		itemToDelete = item
		deleteModal.showModal()
	};

    const confirmDelete = async () => {
        try {
            await pb.collection(itemToDelete.collectionName).delete(itemToDelete.id);
            deleteModal.close();
            fetchMyData(mode);
        } catch (err) {
            alert('Delete failed. Please try again.')
            console.error('Delete error:', err);
        }
    }

	async function fetchMyData(currentMode: ViewMode) {
		// 安全守卫：如果没有登录，直接跳走
		if (!pb.authStore.isValid || !pb.authStore.record) {
			goto('/login');
			return;
		}

		isLoading = true;
		const userId = pb.authStore.record.id;

		try {
			const collection = currentMode === 'cps' ? 'cps' : 'characters';
			const options = {
				// 关键点：只获取 owner 是当前用户的记录
				filter: `owner = "${userId}"`,
				sort: '-created',
				expand: currentMode === 'cps' ? 'characters' : ''
			};

			items = await pb.collection(collection).getFullList(options);
		} catch (err) {
			console.error('Fetch error:', err);
		} finally {
			isLoading = false;
		}
	}

	// 监听 mode 变化自动刷新
	$effect(() => {
		fetchMyData(mode);
	});
</script>

<div class="mx-auto max-w-6xl px-6 py-16">
	<div class="mb-12">
		<h1 class="text-4xl font-bold md:text-5xl">My Stuff</h1>
		<p class="mt-2 text-base-content/60">Manage the CPs and Characters you've created.</p>

		<div role="tablist" class="tabs-boxed mt-6 tabs w-fit">
			<button
				role="tab"
				class="tab {mode === 'cps' ? 'tab-active' : ''}"
				onclick={() => (mode = 'cps')}>My CPs</button
			>
			<button
				role="tab"
				class="tab {mode === 'characters' ? 'tab-active' : ''}"
				onclick={() => (mode = 'characters')}>My Characters</button
			>
		</div>
	</div>

	{#if isLoading}
		<div class="flex justify-center py-20">
			<span class="loading loading-lg loading-spinner text-primary"></span>
		</div>
	{:else if items.length === 0}
		<div class="card border-2 border-dashed border-base-300 bg-base-200 py-20 text-center">
			<div class="card-body items-center">
				<h2 class="card-title text-2xl opacity-40">
					You haven't created any {mode === 'cps' ? 'CP' : 'characters'} yet.
				</h2>
				<div class="mt-4 card-actions">
					<a href="/create" class="btn btn-primary"
						>Create Your First {mode === 'cps' ? 'CP' : 'character'}</a
					>
				</div>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each items as item}
				<div class="card border border-base-200 bg-base-100 shadow-xl">
					<div class="card-body p-5">
						<h2 class="card-title">{item.name}</h2>
						<div class="mt-4 card-actions flex justify-end gap-1">
							<button class="btn w-14 btn-soft btn-sm btn-error" onclick={() => handleDelete(item)}>
								Delete
							</button>

							<a href="/edit/{mode}/{item.id}" class="btn w-14 btn-ghost btn-sm"> Edit </a>

							<a href="/{mode}/{item.id}" class="btn w-14 btn-sm btn-primary"> View </a>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<dialog bind:this={deleteModal} class="modal">
	<div class="modal-box border border-red-200 bg-base-100">
		<h3 class="text-content text-lg font-bold">Confirm Deletion</h3>
		<p class="text-content py-4">
			Are you sure you want to delete this? This action cannot be undone.
		</p>
		{#if itemToDelete?.collectionName! == 'cps'}
			<p class="text-content py-4">
				CP cannot be restored after deletion; the character still exists.
			</p>
		{:else}
			<p class="text-content py-4">
				Once a character is deleted, it cannot be recovered, and the associated CP will also be
				deleted.
			</p>
		{/if}

		<div class="modal-action">
			<form method="dialog">
				<button class="btn btn-ghost">Cancel</button>
			</form>
			<button class="btn text-white btn-soft btn-error" onclick={confirmDelete}> Delete </button>
		</div>
	</div>

	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
