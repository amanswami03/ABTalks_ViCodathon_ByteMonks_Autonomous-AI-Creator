export default function Input({ className = '', ...props }) {
  return (
    <input
      className={`w-full rounded-2xl border border-[#948979]/20 bg-[#222831] px-4 py-3 text-[#DFD0B8] outline-none transition focus:border-[#948979] ${className}`.trim()}
      {...props}
    />
  );
}
