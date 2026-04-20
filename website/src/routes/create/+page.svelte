<script lang="ts">
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import CharacterCreateModal from '$lib/components/CharacterCreateModal.svelte';
	import TagInput from '$lib/components/ui/TagInput.svelte';
	import { toFormData } from '$lib/utils/api';

	let errorText = $state('');
	let name = $state('');
	let description = $state('');
	let pictures: FileList | undefined = $state();
	let character1 = $state('');
	let character2 = $state('');

	let tags = $state<string[]>([]);

	let isNewCharDialogOpen = $state(false);

	const createCP = async () => {
		errorText = '';

		if (
			!name.trim() ||
			!description.trim() ||
			!pictures ||
			pictures.length < 1 ||
			pictures.length > 3 ||
			!character1.trim() ||
			!character2.trim()
		) {
			errorText = 'All fields must be filled out';
			return;
		}

		if (character1 === character2) {
			errorText = 'Character #1 and Character #2 cannot be the same';
			return;
		}
		try {
			await pb
				.collection('cps')
				.create(
					toFormData({
						name,
						description,
						owner: pb.authStore.record!.id,
						images: pictures,
						tag_names: tags,
						characters: [character1, character2]
					})
				);
			window.location.pathname = '/explore';
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Create failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
				return;
			}
		}
	};

	const handleChangeSelect = (e: Event) => {
		const target = e.target as HTMLSelectElement;
		if (target.value === 'new') {
			target.value = '';
			isNewCharDialogOpen = true;
			return;
		}
	};

	// Get Char(s)
	let characters = $state<any[]>([]);
	const reloadCharacters = async () => {
		characters = await pb.collection('characters').getFullList({
			filter: `owner = "${pb.authStore.record!.id}"`,
			sort: '-created'
		});
	};

	onMount(() => {
		if (!pb.authStore.isValid) {
			window.location.pathname = '/login';
			return;
		}
		reloadCharacters();
	});
</script>

<CharacterCreateModal bind:open={isNewCharDialogOpen} afterCreate={() => reloadCharacters()} />

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-12">
	<div class="card w-full max-w-md bg-base-100">
		<div class="card-body p-8">
			<h2 class="mb-2 text-center text-3xl font-bold">Create a CP</h2>
			<p class="mb-6 text-center text-base-content/60">
				All fields must be filled out unless otherwise specified.
			</p>
			<div class="flex flex-col gap-4">
				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your CP Name</span></div>
					<input
						type="text"
						name="cpName"
						placeholder="XXXX & YYYY"
						class="input-bordered input w-full"
						autocomplete="name"
						bind:value={name}
						required
					/>
				</label>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your CP Description</span></div>
					<textarea
						name="cpDescription"
						placeholder="description"
						class="textarea w-full"
						maxlength="1000"
						bind:value={description}
						required
					></textarea>
				</label>

				<div class="form-control w-full max-w-xs">
					<label class="label" for="file-upload">
						<span class="label-text">Your CP Pictures</span>
					</label>

					<input
						type="file"
						id="file-upload"
						class="file-input-bordered file-input w-full file-input-primary"
						accept="image/*"
						bind:files={pictures}
						multiple
						required
					/>
				</div>

				<TagInput title="Your CP Tags" bind:tags />

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your Character #1</span></div>
					<select
						name="cpCharacter1"
						placeholder="Select a character"
						class="select-bordered select w-full"
						autocomplete="name"
						onchange={handleChangeSelect}
						bind:value={character1}
						required
					>
						<option value="" selected disabled>Select a character</option>
						<option value="new">Create a new character</option>
						{#each characters as character}
							<option value={character.id}>{character.name}</option>
						{/each}
					</select>
				</label>

				<label class="form-control w-full">
					<div class="label"><span class="label-text font-medium">Your Character #2</span></div>
					<select
						name="cpCharacter2"
						placeholder="Select a character"
						class="select-bordered select w-full"
						autocomplete="name"
						onchange={handleChangeSelect}
						bind:value={character2}
						required
					>
						<option value="" selected disabled>Select a character</option>
						<option value="new">Create a new character</option>
						{#each characters as character}
							<option value={character.id}>{character.name}</option>
						{/each}
					</select>
				</label>

				<p class="text-red-500">{errorText}</p>

				<div class="form-control mt-6">
					<button onclick={createCP} type="submit" class="btn w-full text-lg btn-primary"
						>Create</button
					>
				</div>
			</div>
		</div>
	</div>
</div>
